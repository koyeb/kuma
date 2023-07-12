import{A as x,_ as b,S as D,a as O}from"./SubscriptionHeader.vue_vue_type_script_setup_true_lang-dea596cb.js";import{a as B,D as C}from"./DefinitionListItem-17aeb6ea.js";import{E as c}from"./EnvoyData-5253b812.js";import{T as S}from"./TabsWidget-2ad3bb32.js";import{T as g}from"./TextWithCopyButton-289007e3.js";import{j as T}from"./RouteView.vue_vue_type_script_setup_true_lang-ed887f62.js";import{d as A,c as m,r as I,o as t,a as l,w as e,k as q,g as d,h as s,t as h,e as v,j as y,b as L,F as p}from"./index-231ca628.js";const V={class:"entity-heading"},H=A({__name:"ZoneEgressDetails",props:{zoneEgressOverview:{type:Object,required:!0}},setup(f){const u=f,{t:_}=T(),z=[{hash:"#overview",title:"Overview"},{hash:"#insights",title:"Zone Egress Insights"},{hash:"#xds-configuration",title:"XDS Configuration"},{hash:"#envoy-stats",title:"Stats"},{hash:"#envoy-clusters",title:"Clusters"}],E=m(()=>({name:"zone-egress-detail-view",params:{zoneEgress:u.zoneEgressOverview.name}})),n=m(()=>{const{type:i,name:a}=u.zoneEgressOverview;return{type:i,name:a}}),k=m(()=>{var a;const i=((a=u.zoneEgressOverview.zoneEgressInsight)==null?void 0:a.subscriptions)??[];return Array.from(i).reverse()});return(i,a)=>{const w=I("router-link");return t(),l(S,{tabs:z},{tabHeader:e(()=>[q("h1",V,[d(`
        Zone Egress:

        `),s(g,{text:n.value.name},{default:e(()=>[s(w,{to:E.value},{default:e(()=>[d(h(n.value.name),1)]),_:1},8,["to"])]),_:1},8,["text"])])]),overview:e(()=>[s(C,null,{default:e(()=>[(t(!0),v(p,null,y(n.value,(o,r)=>(t(),l(B,{key:r,term:L(_)(`http.api.property.${r}`)},{default:e(()=>[r==="name"?(t(),l(g,{key:0,text:o},null,8,["text"])):(t(),v(p,{key:1},[d(h(o),1)],64))]),_:2},1032,["term"]))),128))]),_:1})]),insights:e(()=>[s(O,{"initially-open":0},{default:e(()=>[(t(!0),v(p,null,y(k.value,(o,r)=>(t(),l(x,{key:r},{"accordion-header":e(()=>[s(b,{details:o},null,8,["details"])]),"accordion-content":e(()=>[s(D,{details:o,"is-discovery-subscription":""},null,8,["details"])]),_:2},1024))),128))]),_:1})]),"xds-configuration":e(()=>[s(c,{"data-path":"xds","zone-egress-name":n.value.name,"query-key":"envoy-data-zone-egress"},null,8,["zone-egress-name"])]),"envoy-stats":e(()=>[s(c,{"data-path":"stats","zone-egress-name":n.value.name,"query-key":"envoy-data-zone-egress"},null,8,["zone-egress-name"])]),"envoy-clusters":e(()=>[s(c,{"data-path":"clusters","zone-egress-name":n.value.name,"query-key":"envoy-data-zone-egress"},null,8,["zone-egress-name"])]),_:1})}}});export{H as _};