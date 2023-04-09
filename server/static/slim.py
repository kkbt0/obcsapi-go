
# full 580k lines mem 200-300Mb
# 200k lines mem 120Mb
# 100k lines mem 70Mb
# 20k lines mem 30Mb
# 10k lines mem 27Mb
with open('dictionary.txt', 'r') as input_file, \
     open('dictionary-20k.txt', 'w') as output_file:
    sum = 0
    for line in input_file:
        if sum <= 20000:
            output_file.write(line)
            sum += 1

# 汉字 [\u3400-\u4db5]
# [\u3400-\u4db5] 汉字扩充（繁体字、不常见字） 
# 常规汉字 [\u4e00-\u9fa5]