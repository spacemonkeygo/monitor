#include "runtime.h"

void ·schedTrace(Slice b, bool detailed, int32 *n) {
  if(b.len == 0) {
    *n = 0;
    return;
  }

  g->writebuf = (byte*)b.array;
  g->writenbuf = b.len;
  runtime·schedtrace(detailed);
  *n = b.len - g->writenbuf;
  g->writebuf = nil;
  g->writenbuf = 0;
}
