const e="HealthCheck",t="default",s="web-to-banana",n=[{match:{service:"web"}}],a=[{match:{service:"backend"}}],c={activeChecks:{interval:"10s",timeout:"2s",unhealthyThreshold:123,healthyThreshold:12}},o={type:e,mesh:t,name:s,sources:n,destinations:a,conf:c};export{c as conf,o as default,a as destinations,t as mesh,s as name,n as sources,e as type};
