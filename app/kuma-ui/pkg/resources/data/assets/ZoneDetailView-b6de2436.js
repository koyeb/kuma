import{d as _,u as d,r as i,v as u,o,j as c,b as l,g as k}from"./index-8d80a271.js";import{_ as w}from"./ZoneDetails.vue_vue_type_script_setup_true_lang-69838ecd.js";import{_ as z}from"./EmptyBlock.vue_vue_type_script_setup_true_lang-92154385.js";import{E as h}from"./ErrorBlock-755be44c.js";import{_ as y}from"./LoadingBlock.vue_vue_type_script_setup_true_lang-b1bc97a7.js";import{u as g}from"./store-5737c223.js";import{u as B}from"./index-4917657d.js";import"./kongponents.es-176426e1.js";import"./AccordionList-c29101b8.js";import"./_plugin-vue_export-helper-c27b6911.js";import"./CodeBlock.vue_vue_type_style_index_0_lang-aa642a42.js";import"./DefinitionListItem-56ada3b3.js";import"./SubscriptionHeader.vue_vue_type_script_setup_true_lang-478c9728.js";import"./TabsWidget-1dcc8ca1.js";import"./datadogLogEvents-302eea7b.js";import"./QueryParameter-70743f73.js";import"./TextWithCopyButton-40f68948.js";import"./WarningsWidget.vue_vue_type_script_setup_true_lang-df13d6f4.js";const E={class:"zone-details"},$={key:3,class:"kcard-border"},H=_({__name:"ZoneDetailView",setup(b){const p=B(),e=d(),f=g(),r=i(null),n=i(!0),a=i(null);u(()=>e.params.mesh,function(){e.name==="zone-detail-view"&&s()}),u(()=>e.params.name,function(){e.name==="zone-detail-view"&&s()}),v();function v(){f.dispatch("updatePageTitle",e.params.zone),s()}async function s(){n.value=!0,a.value=null;const m=e.params.zone;try{r.value=await p.getZoneOverview({name:m})}catch(t){r.value=null,t instanceof Error?a.value=t:console.error(t)}finally{n.value=!1}}return(m,t)=>(o(),c("div",E,[n.value?(o(),l(y,{key:0})):a.value!==null?(o(),l(h,{key:1,error:a.value},null,8,["error"])):r.value===null?(o(),l(z,{key:2})):(o(),c("div",$,[k(w,{"zone-overview":r.value},null,8,["zone-overview"])]))]))}});export{H as default};