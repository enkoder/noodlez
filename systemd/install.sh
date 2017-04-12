sudo mkdir -p /etc/systemd/system/
sudo cp fcserver.service /etc/systemd/system/
sudo cp noodle.service /etc/systemd/system/
sudo systemctl enable noodle.service
sudo systemctl enable fcserver.service
