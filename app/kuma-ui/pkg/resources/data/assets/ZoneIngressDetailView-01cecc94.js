import{d as h,a as m,o,b as c,w as t,e as a,p as u,a1 as g,f as n,t as i,q as v,c as l,F as p,C as z,v as x}from"./index-784d2bbf.js";import{S as C}from"./StatusBadge-a6acfbee.js";import{T as w}from"./TextWithCopyButton-7ef74197.js";import{_ as V}from"./SubscriptionList.vue_vue_type_script_setup_true_lang-45014a53.js";import{g as b}from"./dataplane-dcd0858b.js";import"./CopyButton-9c00109a.js";import"./index-9dd3e7d3.js";import"./AccordionList-e1625b82.js";const B={class:"stack","data-testid":"detail-view-details"},A={class:"columns"},$={key:0},L=h({__name:"ZoneIngressDetailView",props:{data:{}},setup(y){const e=y;return(N,S)=>{const _=m("KCard"),I=m("AppView"),f=m("RouteView");return o(),c(f,{name:"zone-ingress-detail-view"},{default:t(({t:r})=>[a(I,null,{default:t(()=>{var k;return[u("div",B,[a(_,null,{body:t(()=>[u("div",A,[a(g,null,{title:t(()=>[n(i(r("http.api.property.status")),1)]),body:t(()=>[a(C,{status:v(b)(e.data.zoneIngressInsight)},null,8,["status"])]),_:2},1024),n(),a(g,null,{title:t(()=>[n(i(r("http.api.property.address")),1)]),body:t(()=>{var s,d;return[(s=e.data.zoneIngress.networking)!=null&&s.address&&((d=e.data.zoneIngress.networking)!=null&&d.port)?(o(),c(w,{key:0,text:`${e.data.zoneIngress.networking.address}:${e.data.zoneIngress.networking.port}`},null,8,["text"])):(o(),l(p,{key:1},[n(i(r("common.detail.none")),1)],64))]}),_:2},1024),n(),a(g,null,{title:t(()=>[n(i(r("http.api.property.advertisedAddress")),1)]),body:t(()=>{var s,d;return[(s=e.data.zoneIngress.networking)!=null&&s.advertisedAddress&&((d=e.data.zoneIngress.networking)!=null&&d.advertisedPort)?(o(),c(w,{key:0,text:`${e.data.zoneIngress.networking.advertisedAddress}:${e.data.zoneIngress.networking.advertisedPort}`},null,8,["text"])):(o(),l(p,{key:1},[n(i(r("common.detail.none")),1)],64))]}),_:2},1024)])]),_:2},1024),n(),(o(!0),l(p,null,z([((k=e.data.zoneIngressInsight)==null?void 0:k.subscriptions)??[]],s=>(o(),l(p,{key:s},[s.length>0?(o(),l("div",$,[u("h2",null,i(r("zone-ingresses.routes.item.subscriptions.title")),1),n(),a(_,{class:"mt-4"},{body:t(()=>[a(V,{subscriptions:s},null,8,["subscriptions"])]),_:2},1024)])):x("",!0)],64))),128))])]}),_:2},1024)]),_:1})}}});export{L as default};