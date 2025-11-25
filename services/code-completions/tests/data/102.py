import subprocess
import re

def run_command(command):
    """运行命令并返回输出"""
    try:
        result = subprocess.run(command, stdout=subprocess.PIPE, stderr=subprocess.PIPE, text=True, shell=True)
        if result.returncode != 0:
            print(f"命令执行失败: {result.stderr}")
            return None
        return result.stdout
    except Exception as e:
        print(f"运行命令时发生错误: {e}")
        return None

def get_images_with_pattern(pattern):
    """获取包含指定模式的镜像"""
    output = run_command("docker images")
    if not output:
        return []
    
    images = []
    for line in output.splitlines():
        if re.search(pattern, line):
            images.append(line)
    return images

def main():
    pattern = "zgsm/code-completion"
    images = get_images_with_pattern(pattern)
    for image in images:
        print(image)

if __name__ == "__main__":
    
