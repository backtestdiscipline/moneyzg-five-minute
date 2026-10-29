import csv,zipfile,json,hashlib,pathlib,math,collections,sys
root=pathlib.Path.cwd();data=root/'data/development';manifest=json.loads((root/'provenance/development.json').read_bytes());first=1704067200000;end=1735516800000;warm=first-500*300000
bars=[];fund=[]
for m in manifest:
 p=data/m['name']; raw=p.read_bytes(); assert hashlib.sha256(raw).hexdigest()==m['sha256'] and len(raw)==m['bytes']
 if m['name'].endswith('.zip'):
  with zipfile.ZipFile(p) as z:
   with z.open(z.namelist()[0]) as f:
    import io
    for row in csv.reader(io.TextIOWrapper(f)):
     if row[0]=='open_time':continue
     t=int(row[0]);
     if warm<=t<end:bars.append((t,*map(float,row[1:6])))
 elif m['name'].startswith('funding-page-'):
  for x in json.loads(raw): fund.append((int(x['fundingTime']),float(x['fundingRate']),float(x['markPrice'])))
bars.sort();assert len(bars)==105332 and all(bars[k][0]==warm+k*300000 for k in range(len(bars)))
# Build complete causal signal event list first: all inputs in each signal window end at confirmation.
e20=[];e50=[]
for i,bar in enumerate(bars):
 c=bar[4]
 e20.append(c if i==0 else e20[-1]+(2/21)*(c-e20[-1]));e50.append(c if i==0 else e50[-1]+(2/51)*(c-e50[-1]))
signals={}
for p in range(500,len(bars)-3):
 lo=bars[p][3];hi=bars[p][2];neighbors=[bars[j] for j in (p-2,p-1,p+1,p+2)]
 islow=all(lo<b[3] for b in neighbors);ishigh=all(hi>b[2] for b in neighbors);c=p+2
 if islow and not ishigh and e50[p]<lo<e20[p] and e20[p]>e50[p] and e20[c]>e50[c]:signals[c]=(1,e50[c],p)
 if ishigh and not islow and e20[p]<hi<e50[p] and e20[p]<e50[p] and e20[c]<e50[c]:signals[c]=(-1,e50[c],p)
funds={x[0]-x[0]%28800000:x for x in fund};trades=[];position=None;pending=None;ignored=invalid=ambig_count=0;funding_events=[]
def complete(p,bar,raw,reason,amb):
 direction=p['Direction']; exitprice=raw*(1-direction*0.0002);fee=(p['Entry']+exitprice)*0.0005;slip=(p['raw_entry']+raw)*0.0002
 p.update(ExitTime=bar[0],Exit=exitprice,GrossR=direction*(raw-p['raw_entry'])/p['Risk'],NetR=(direction*(exitprice-p['Entry'])-fee-p['FundingCost'])/p['Risk'],Fees=fee,SlippageCost=slip,Ambiguous=amb,ExitReason=reason)
 return {k:v for k,v in p.items() if k!='raw_entry'}
def open_exit(p,b):
 d=p['Direction'];o=b[1]
 if (o-p['Stop'])*d<=0:return o,'stop',False
 if (o-p['Target'])*d>=0:return p['Target'],'target',False
 return None
def interior_exit(p,b):
 d=p['Direction']; stop=(b[3]<=p['Stop'] if d==1 else b[2]>=p['Stop']); target=(b[2]>=p['Target'] if d==1 else b[3]<=p['Target'])
 if stop:return p['Stop'],'stop',target
 if target:return p['Target'],'target',False
 return None
for i in range(500,len(bars)):
 b=bars[i]; f=funds.get(b[0])
 if position and f and f[0]==b[0]:position['FundingCost']+=position['Direction']*f[1]*f[2]
 if position:
  x=open_exit(position,b)
  if x:trades.append(complete(position,b,*x));position=None
 if pending:
  d,stop,p=pending;entry=b[1]*(1+d*0.0002);risk=d*(entry-stop);target=entry+d*1.5*risk
  if risk<=0 or target<=0:invalid+=1
  else:position=dict(Direction=d,EntryTime=b[0],Entry=entry,Stop=stop,Target=target,Risk=risk,FundingCost=0,FundingAmbiguous=False,raw_entry=b[1])
  pending=None
 if position:
  x=open_exit(position,b)
  if x:trades.append(complete(position,b,*x));position=None
 if position:
  x=interior_exit(position,b)
  if f and f[0]>b[0]:
   c=position['Direction']*f[1]*f[2]
   if x:position['FundingAmbiguous']=True;position['FundingCost']+=max(c,0)
   else:position['FundingCost']+=c
  if x:
   trades.append(complete(position,b,*x));ambig_count+=int(x[2]);position=None
 if i+1<len(bars) and i in signals:
  if position:ignored+=1
  else:pending=signals[i]
out=dict(Trades=trades,InvalidEntries=invalid,IgnoredSignals=ignored,AmbiguousTrades=ambig_count,OpenPosition=position is not None)
ref=json.load(open(sys.argv[1]));print('counts independent',len(trades),invalid,ignored,ambig_count,position is not None);print('counts reference',len(ref['Trades']),ref['InvalidEntries'],ref['IgnoredSignals'],ref['AmbiguousTrades'],ref['OpenPosition']);assert all(out[k]==ref[k] for k in out if k!='Trades'); assert len(trades)==len(ref['Trades'])
maxdiff=collections.defaultdict(float);problems=[]
for n,(a,b) in enumerate(zip(trades,ref['Trades'])):
 for k,v in b.items():
  if isinstance(v,float):
   diff=abs(a[k]-v);maxdiff[k]=max(maxdiff[k],diff)
   if diff>1e-7*max(1,abs(v)):problems.append((n,k,a[k],v))
  elif a[k]!=v:problems.append((n,k,a[k],v))
print('max_abs_float_differences',dict(maxdiff));print('mismatches',len(problems),'sample',problems[:5]);assert not problems
