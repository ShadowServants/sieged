from collections import defaultdict


def hemming_rast(a, b):
    diff = 0
    for i in range(len(a)):
        if a[i] != b[i]:
            diff += 1
    return diff


def padding(k):
    return '0' * (9 - len(k)) + k


'2 есть'
b = 'xxx**xxxx'
a = '001101101'
b = '001011001'

a_dict = defaultdict(int)

for i in range(0, 2 ** 9 - 1):
    numm = padding(bin(i)[2:])

    a_dict[hemming_rast(numm, a)] += 1

b_dict = defaultdict(int)

for i in range(0, 2 ** 9 - 1):
    numm = padding(bin(i)[2:])

    b_dict[hemming_rast(numm, b)] += 1


print(a_dict)
print(b_dict)
summ = 0
for i in range(0, 8):
    summ += min(a_dict[i], b_dict[8 - i])

print(summ)
