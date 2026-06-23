import subprocess

subprocess.run(["cmake", "-B", "build"])
subprocess.run(["cmake", "--build", "build"])
