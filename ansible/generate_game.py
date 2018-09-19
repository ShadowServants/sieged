import os
import shutil

import click
from jinja2 import Environment, PackageLoader

env = Environment(
    loader=PackageLoader('generate_game', 'templates')
)

services = []


def create_dir_if_not_exists(path):
    if not os.path.exists(path):
        os.makedirs(path)


def copy_file(from_path, to_path):
    shutil.copyfile(from_path, to_path)


def copy_binary(from_path, to_path):
    shutil.copyfile(from_path, to_path)
    os.chmod(to_path, 0o755)


REDIS_PORT = 6380
FLAG_HANDLER_PORT = 7000
FLAG_ADDER_PORT = 5000
ROUND_HANDLER_PORT = 8050

binaries = ['flag_adder', 'flag_handler', 'round_handler']


def get_defaults_dict(*args, **kwargs):
    return dict(redis_host=kwargs.get("redis_host") or "127.0.0.1",
                redis_port=kwargs.get("redis_port") or "6379",
                team_num=kwargs.get("team_num") or 20,
                flag_adder_host="127.0.0.1",
                round_handler_host="127.0.0.1")


def render_template(template_name, *args, **kwargs):
    template = env.get_template(template_name)
    config = get_defaults_dict(*args, **kwargs)
    config.update(kwargs)
    return template.render(**config)


def save_config_file(filename, content):
    with open(filename, 'w') as file:
        file.write(content)


@click.command()
@click.option('--path', default='supervisor', help='Folder to create applications')
@click.option('--platform', default='darwin64', help='Platform to run available: darwin64/amd64')
def create_service(path, platform):
    global FLAG_ADDER_PORT, FLAG_HANDLER_PORT, ROUND_HANDLER_PORT, REDIS_PORT
    name = click.prompt("Enter service name")

    click.echo("Creating service-controller directory for {}".format(name))
    service_path = path + "/service_controller-{}".format(name)
    create_dir_if_not_exists(service_path)
    click.echo("Created directory {}".format(service_path))
    click.echo("Copy binaries for service")
    with click.progressbar(binaries) as progress:
        for binary in progress:
            from_path = "supervisor/binaries/{}/{}".format(platform, binary)
            to_path = service_path + "/{}".format(binary)
            copy_binary(from_path, to_path)
            click.echo(" Copy {} to {}".format(from_path, to_path))

    click.echo("Create symlink for teams file")
    os.symlink("../teams.yaml", service_path + "/teams.yaml")

    # Configure redis

    REDIS_PORT = int(click.prompt("Enter port for redis", default=str(REDIS_PORT)))
    redis_conf = render_template("redis.cnf", redis_port=REDIS_PORT,service_name=name)
    save_config_file(service_path + "/" + "redis.cnf", redis_conf)

    service_controller_conf = render_template("service_controller.conf", redis_port=REDIS_PORT,
                                              service_directory="service_controller-{}".format(name),
                                              service_name="service_controller-{}".format(name))
    save_config_file(service_path + "/" + "service_controller.conf", service_controller_conf)

    # Configure flag adder
    click.echo("Configuring flag adder")
    flag_prefix = click.prompt("Enter flag prefix", default=name.upper()[0])
    FLAG_ADDER_PORT = int(click.prompt("Enter port for flag adder", default=str(FLAG_ADDER_PORT)))
    flag_adder_conf = render_template("flag_adder.yaml", flag_adder_port=FLAG_ADDER_PORT, redis_port=REDIS_PORT,
                                      flag_prefix=flag_prefix)

    save_config_file(service_path + "/" + "flag_adder.yaml", flag_adder_conf)
    # Configure round_handler
    click.echo("Configuring round handler")
    ROUND_HANDLER_PORT = int(click.prompt("Enter port for round handler", default=str(ROUND_HANDLER_PORT)))
    checker_name = click.prompt("Enter checker name", default=name + "_check.py")
    round_handler_conf = render_template("round_handler.yaml", round_handler_port=ROUND_HANDLER_PORT,
                                         redis_port=REDIS_PORT, checker_name=checker_name)
    save_config_file(service_path + "/" + "round_handler.yaml", round_handler_conf)
    # Configure flag_handler
    click.echo("Configuring flag handler")
    FLAG_HANDLER_PORT = int(click.prompt("Enter port for round handler", default=str(FLAG_HANDLER_PORT)))
    flag_handler_conf = render_template("flag_handler.yaml", flag_handler_port=FLAG_HANDLER_PORT,
                                        redis_port=REDIS_PORT)
    save_config_file(service_path + "/" + "flag_handler.yaml", flag_handler_conf)
    services.append(
        dict(prefix=flag_prefix, round_handler_port=ROUND_HANDLER_PORT, flag_handler_port=FLAG_HANDLER_PORT, name=name))
    FLAG_HANDLER_PORT += 1
    FLAG_ADDER_PORT += 1
    REDIS_PORT += 1
    ROUND_HANDLER_PORT += 1


@click.command()
@click.option('--path', default='supervisor', help='Folder to create applications')
@click.option('--num', help='Enter number of services please')
@click.option('--platform', default='darwin64', help='Platform to run available: darwin64/amd64')
@click.pass_context
def generate_game(ctx, path, num, platform):
    create_dir_if_not_exists(path)
    click.echo("Created directory {}".format(path))
    create_dir_if_not_exists(path + "/redis_store")
    click.echo("Created directory {}".format(path + "/redis_store"))
    click.echo("Created supervisord config")
    copy_file("templates/supervisord.conf", path + "/supervisord.conf")
    click.echo("Created teams file")
    copy_file("templates/teams.yaml", path + "/teams.yaml")
    click.echo("Copy binaries for router and tokens")
    with click.progressbar(["http_router", "tokens"]) as progress:
        for binary in progress:
            from_path = "supervisor/binaries/{}/{}".format(platform, binary)
            to_path = path + "/{}".format(binary)
            copy_binary(from_path, to_path)
            click.echo(" Copy {} to {}".format(from_path, to_path))

    for i in range(int(num)):
        click.echo("Creating service #{}".format(str(i + 1)))
        ctx.invoke(create_service, platform=platform, path=path)
    router_conf = render_template("router_config.yaml", router_host="0.0.0.0", router_port="31337", services=services)
    save_config_file(path + "/" + "router_config.yaml", router_conf)
    tasks_conf = render_template("tasks.yaml", services=services)
    save_config_file(path + "/" + "tasks.yaml", tasks_conf)



if __name__ == '__main__':
    generate_game()
