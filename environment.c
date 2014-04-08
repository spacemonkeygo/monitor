#include "runtime.h"

// IStats
// shared with Go. If you edit this struct, edit InternalStats
struct IStats
{
  int32 GoMaxProcs;
  int32 IdleProcs;
  int32 ThreadCount;
  int32 IdleThreads;
  int32 RunQueue;
  int32 ProcRunQueueSize;
  int32 ProcRunQueueTotal;
};

typedef struct IStats IStats;

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

void ·partialRuntimeInternals(IStats *res) {
  int32 i, q, t, h, s, proc_q, proc_t;
  P *p;

  proc_q = 0;
  proc_t = 0;

  for(i = 0; i < runtime·gomaxprocs; i++) {
    p = runtime·allp[i];
    if(p == nil)
      continue;
    t = p->runqtail;
    h = p->runqhead;
    s = p->runqsize;
    q = t - h;
    if(q < 0)
      q += s;
    proc_q += q;
    proc_t += s;
  }

  res->GoMaxProcs = runtime·gomaxprocs;
  res->ThreadCount = runtime·mcount();
  res->ProcRunQueueSize = proc_q;
  res->ProcRunQueueTotal = proc_t;

  // maybe someday we can set the rest without calling schedTrace
}
