import os

from projects_to_files import projects_to_files


COMMITS_PATH="/raid/CNVRG_MINIO_BACKUP2/cnvrg-storage"
ROOT_PATH = "/raid/projects_files"

os.system("mkdir -p " + ROOT_PATH)


def cp_file(file_name, project_name):
    path = file_name.split("/")
    new_path = "/".join(path[4:])
    subdirectories = "/".join(new_path.split("/")[:-1])
    if subdirectories != "":
        os.system("mkdir -p " + ROOT_PATH + "/" + project +  "/" + subdirectories)
    code = os.system("cp '" + COMMITS_PATH + "/" + file_name + "' '" + ROOT_PATH + "/" + project +  "/" + new_path + "'")
    if code != 0:
        print(project, path, new_path)


for project, files in projects_to_files.items():
    os.system("mkdir -p " + ROOT_PATH + "/" + project)
    for file_ in files:
        cp_file(file_, project)


