import{A as _,a as z}from"./AccordionList-1e0cd60d.js";import{D as E,a as b}from"./DefinitionListItem-87124c14.js";import{E as l}from"./EnvoyData-3b512034.js";import{_ as w,S as D}from"./SubscriptionHeader.vue_vue_type_script_setup_true_lang-dceeb786.js";import{T as k}from"./TabsWidget-597e6684.js";import{d as x,c as d,b as c,w as e,o as n,i as O,t as m,g as s,j as v,q as g,h as S,F as p}from"./index-0fbacd76.js";const q={class:"entity-heading"},Z=x({__name:"ZoneEgressDetails",props:{zoneEgressOverview:{type:Object,required:!0}},setup(h){const u=h,y=[{hash:"#overview",title:"Overview"},{hash:"#insights",title:"Zone Egress Insights"},{hash:"#xds-configuration",title:"XDS Configuration"},{hash:"#envoy-stats",title:"Stats"},{hash:"#envoy-clusters",title:"Clusters"}],t=d(()=>{const{type:r,name:a}=u.zoneEgressOverview;return{type:r,name:a}}),f=d(()=>{var a;const r=((a=u.zoneEgressOverview.zoneEgressInsight)==null?void 0:a.subscriptions)??[];return Array.from(r).reverse()});return(r,a)=>(n(),c(k,{tabs:y},{tabHeader:e(()=>[O("h1",q,`
        Zone Egress: `+m(t.value.name),1)]),overview:e(()=>[s(b,null,{default:e(()=>[(n(!0),v(p,null,g(t.value,(o,i)=>(n(),c(E,{key:i,term:i},{default:e(()=>[S(m(o),1)]),_:2},1032,["term"]))),128))]),_:1})]),insights:e(()=>[s(z,{"initially-open":0},{default:e(()=>[(n(!0),v(p,null,g(f.value,(o,i)=>(n(),c(_,{key:i},{"accordion-header":e(()=>[s(w,{details:o},null,8,["details"])]),"accordion-content":e(()=>[s(D,{details:o,"is-discovery-subscription":""},null,8,["details"])]),_:2},1024))),128))]),_:1})]),"xds-configuration":e(()=>[s(l,{"data-path":"xds","zone-egress-name":t.value.name,"query-key":"envoy-data-zone-egress"},null,8,["zone-egress-name"])]),"envoy-stats":e(()=>[s(l,{"data-path":"stats","zone-egress-name":t.value.name,"query-key":"envoy-data-zone-egress"},null,8,["zone-egress-name"])]),"envoy-clusters":e(()=>[s(l,{"data-path":"clusters","zone-egress-name":t.value.name,"query-key":"envoy-data-zone-egress"},null,8,["zone-egress-name"])]),_:1}))}});export{Z as _};
