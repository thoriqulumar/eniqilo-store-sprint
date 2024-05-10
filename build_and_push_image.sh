# check if all arguments passed to the script
if [ $# -ne 2 ]; then
    echo "Usage: $0 <your_username> <your_token>"
    exit 1
fi

USERNAME=$1
PASSWORD=$2

docker build -t main_thorumr:latest .

docker login -u $USERNAME -p $PASSWORD docker.io

docker tag main_thorumr docker.io/$USERNAME/main_thorumr:latest

docker push docker.io/$USERNAME/main_thorumr:latest

