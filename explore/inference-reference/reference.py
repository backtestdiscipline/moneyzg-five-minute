import math, random

def interval(s,n):
    if len(s)!=365 or len(n)!=365 or any(x<0 for x in n): raise ValueError('bad inputs')
    total=sum(n)
    if not total: raise ValueError('no trades')
    theta=sum(s)/total
    z=[sum(s[i:i+73])-theta*sum(n[i:i+73]) for i in range(0,365,73)]
    se=math.sqrt(1.25*sum(x*x for x in z))/total
    return theta-2.7764451051977987*se,theta+2.7764451051977987*se

def wilson(k,n):
    z=1.959963984540054;p=k/n;den=1+z*z/n;mid=(p+z*z/(2*n))/den;delta=z*math.sqrt(p*(1-p)/n+z*z/(4*n*n))/den
    return mid-delta,mid+delta
s=[float(i) for i in range(1,6) for _ in range(73)]; n=[1]*365
lo,hi=interval(s,n);want=2.7764451051977987*math.sqrt(12.5)/5
assert abs(lo-(3-want))<1e-12 and abs(hi-(3+want))<1e-12

cases=('iid','t3','ar05','ar08','variable','cluster')
for case in cases:
    rng=random.Random(961702);reps=3000;hit=0;invalid=0
    for rep in range(reps):
        shocks=[];counts=[];x=rng.gauss(0,1);vol=rng.gauss(0,1)
        for day in range(365):
            if case=='t3':
                g=rng.gauss(0,1);chi=sum(rng.gauss(0,1)**2 for _ in range(3));x=g/math.sqrt(chi) # t3 scaled to unit variance
            elif case in ('ar05','ar08','variable'):
                phi=0.5 if case=='ar05' else 0.8
                x=phi*x+math.sqrt(1-phi*phi)*rng.gauss(0,1)
            elif case=='cluster':
                vol=.9*vol+math.sqrt(1-.9**2)*rng.gauss(0,1)
                x=math.exp(.5*vol-.25)*rng.gauss(0,1) # E sigma^2 = 1
            else:x=rng.gauss(0,1)
            k=rng.randrange(5) if case=='variable' else 1
            shocks.append(k*x);counts.append(k)
        lo,hi=interval(shocks,counts)
        hit+=lo<=0<=hi
    a,b=wilson(hit,reps)
    print(case,hit,reps,f'{hit/reps:.6f}',f'wilson95=[{a:.6f},{b:.6f}]',flush=True)
