wireless_interface=$(ls /sys/class/net/ | grep '^wlan' | head -n 1)
if [ -n "$wireless_interface" ]; then
    mac_address=$(cat /sys/class/net/"$wireless_interface"/address)
    echo "SUBSYSTEM==\"net\", ACTION==\"add\", ATTR{address}==\"$mac_address\", NAME=\"wlan0\"" > /etc/udev/rules.d/70-wlan0-mac.rules
    sudo udevadm control --reload
    echo "Mapped MAC address $mac_address to wlan0."
else
    echo "No wireless interface starting with 'wlan' found."
fi