import{d as y,O as l,k as f,f as s,a as d,H as v,n as r,w as c,c as E,p as i,o as a,g as o,t,F as w,i as B,b as u,D as b,N as x,Q as g,y as S,z as I,e as N}from"./index.d0c2efa7.js";const h=e=>(S("data-v-6e270615"),e=e(),I(),e),C={class:"error-block"},V={class:"card-icon mb-3"},D=h(()=>o("p",null,"An error has occurred while trying to load this data.",-1)),K=h(()=>o("summary",null,"Details",-1)),z={key:0},A={key:0,class:"badge-list"},F=y({__name:"ErrorBlock",props:{error:{type:[Error,l],required:!1,default:null}},setup(e){const n=e,_=f(()=>n.error instanceof Error),p=f(()=>n.error instanceof l?n.error.causes:[]);return(O,j)=>(a(),s("div",C,[d(r(x),{"cta-is-hidden":""},v({title:c(()=>[o("div",V,[d(r(b),{class:"kong-icon--centered",icon:"warning",color:"var(--black-75)","secondary-color":"var(--yellow-300)",size:"42"})]),D]),_:2},[r(_)||r(p).length>0?{name:"message",fn:c(()=>[o("details",null,[K,r(_)?(a(),s("p",z,t(e.error.message),1)):i("",!0),o("ul",null,[(a(!0),s(w,null,B(r(p),(m,k)=>(a(),s("li",{key:k},[o("b",null,[o("code",null,t(m.field),1)]),u(": "+t(m.message),1)]))),128))])])]),key:"0"}:void 0]),1024),e.error instanceof r(l)?(a(),s("div",A,[e.error.code?(a(),E(r(g),{key:0,appearance:"warning"},{default:c(()=>[u(t(e.error.code),1)]),_:1})):i("",!0),d(r(g),{appearance:"warning"},{default:c(()=>[u(t(e.error.statusCode),1)]),_:1})])):i("",!0)]))}});const H=N(F,[["__scopeId","data-v-6e270615"]]);export{H as E};
