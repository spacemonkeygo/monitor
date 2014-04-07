#include "runtime.h"

// IStats
// shared with Go. If you edit this struct, edit InternalStats
struct IStats
{
  int32 GoMaxProcs;
  int32 ThreadCount;
  int32 ProcRunQueueSize;
  int32 ProcRunQueueTotal;
};

typedef struct IStats IStats;

void ·runtimeInternals(IStats *res) {
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
}
