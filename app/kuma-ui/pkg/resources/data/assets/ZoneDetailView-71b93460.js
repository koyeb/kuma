import{d as k,u as w,q as n,o as e,a as r,w as c,h as i,b as m,g as z,k as b,e as h}from"./index-bec4d193.js";import{_ as y}from"./ZoneDetails.vue_vue_type_script_setup_true_lang-dc09a161.js";import{h as $,i as x,f as B,_ as E}from"./RouteView.vue_vue_type_script_setup_true_lang-19528cd9.js";import{_ as V}from"./RouteTitle.vue_vue_type_script_setup_true_lang-2f1aad7c.js";import{_ as g}from"./EmptyBlock.vue_vue_type_script_setup_true_lang-ab21e287.js";import{E as N}from"./ErrorBlock-828beaca.js";import{_ as A}from"./LoadingBlock.vue_vue_type_script_setup_true_lang-91d5cd68.js";import"./kongponents.es-b2b550a8.js";import"./SubscriptionHeader.vue_vue_type_script_setup_true_lang-57b622e7.js";import"./DefinitionListItem-d3bb1a9c.js";import"./CodeBlock.vue_vue_type_style_index_0_lang-94e75ba9.js";import"./TabsWidget-64af74c5.js";import"./QueryParameter-70743f73.js";import"./TextWithCopyButton-5e48fe52.js";import"./WarningsWidget.vue_vue_type_script_setup_true_lang-5fd41e98.js";const C={class:"zone-details"},D={key:3,class:"kcard-border","data-testid":"detail-view-details"},S=k({__name:"ZoneDetailView",setup(O){const _=$(),p=w(),{t:l}=x(),a=n(null),s=n(!0),t=n(null);f();function f(){d()}async function d(){s.value=!0,t.value=null;const u=p.params.zone;try{a.value=await _.getZoneOverview({name:u})}catch(o){a.value=null,o instanceof Error?t.value=o:console.error(o)}finally{s.value=!1}}return(u,o)=>(e(),r(E,null,{default:c(({route:v})=>[i(V,{title:m(l)("zone-cps.routes.item.title",{name:v.params.zone})},null,8,["title"]),z(),i(B,{breadcrumbs:[{to:{name:"zone-cp-list-view"},text:m(l)("zone-cps.routes.item.breadcrumbs")}]},{default:c(()=>[b("div",C,[s.value?(e(),r(A,{key:0})):t.value!==null?(e(),r(N,{key:1,error:t.value},null,8,["error"])):a.value===null?(e(),r(g,{key:2})):(e(),h("div",D,[i(y,{"zone-overview":a.value},null,8,["zone-overview"])]))])]),_:1},8,["breadcrumbs"])]),_:1}))}});export{S as default};