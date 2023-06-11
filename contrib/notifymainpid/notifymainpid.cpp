#include <systemd/sd-daemon.h>
#include <format>
#include <iostream>
#include <fstream>

int main(int argc, char *argv[]) {
  if (argc != 2) {
    fprintf(stderr, "error: incorrect number of arguments\n");
    return 1;
  }  
  std::ifstream mainpid_stream(argv[1]);
  pid_t mainpid;
  mainpid_stream >> mainpid;
  char *podmanpidstr = getenv("SYSTEMD_EXEC_PID");
  pid_t podmanpid = atoi(podmanpidstr);
  std::string msg = std::format("MAINPID={}\nREADY=1", mainpid);
  sd_pid_notify(podmanpid, 0, msg.c_str());
  return 0;
}
