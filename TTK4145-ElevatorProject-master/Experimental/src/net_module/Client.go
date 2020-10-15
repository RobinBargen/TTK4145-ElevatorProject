package net

import (
  "./network/localip"
)

func setUpLocalIP() {
  _, err := localip.LocalIP()
  if err != nil {
    // Error
  }
}
