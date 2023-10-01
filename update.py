import os
import fileinput
import shutil

# 假设这是你的字符串列表，每一个字符串都是一个文件路径
file_paths = [
    'step2.sh', 
    'docs/md/go-version/2-运行与部署.md', 
    'obcsapi-web/src/stores/setting.ts',
    'server/assest.go',
    'server/Dockerfile',
    'server/config.example.yaml',
    'server/server.go',
    'server/tools/config.example.yaml',
    ]

for file_path in file_paths:
    # 检查文件是否存在
    if not os.path.isfile(file_path):
        print(f"File {file_path} does not exist.")
        continue

    with fileinput.FileInput(file_path, inplace=True) as file:
        for line in file:
            # 读取每一行并替换其中的 '4.2.7' 为 '4.2.8'
            print(line.replace('4.2.7', '4.2.8'), end='')

shutil.copy('server/config.example.yaml', 'server/tools/config.example.yaml')
shutil.copy('server/config.example.yaml', 'docs/md/go-version/config.example.yaml')