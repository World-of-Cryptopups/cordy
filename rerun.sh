go build .
sudo systemctl restart cordy_bot.service
journalctl -f -u cordy_bot.service