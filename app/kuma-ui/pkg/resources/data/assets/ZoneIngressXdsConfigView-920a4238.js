import{E as d}from"./EnvoyData-cf2afcb8.js";import{d as l,a as n,o as m,b as g,w as t,e as a,p as u,f as _}from"./index-8567ed34.js";import"./index-52545d1d.js";import"./CodeBlock.vue_vue_type_style_index_0_lang-93993f58.js";import"./EmptyBlock.vue_vue_type_script_setup_true_lang-c8b34455.js";import"./ErrorBlock-3040d559.js";import"./TextWithCopyButton-cadc290c.js";import"./CopyButton-faeef64d.js";import"./WarningIcon.vue_vue_type_script_setup_true_lang-184215be.js";import"./LoadingBlock.vue_vue_type_script_setup_true_lang-bbe0c89a.js";const B=l({__name:"ZoneIngressXdsConfigView",setup(f){return(x,h)=>{const s=n("RouteTitle"),r=n("KCard"),i=n("AppView"),p=n("RouteView");return m(),g(p,{name:"zone-ingress-xds-config-view",params:{zoneIngress:"",codeSearch:"",codeFilter:!1,codeRegExp:!1}},{default:t(({route:e,t:c})=>[a(i,null,{title:t(()=>[u("h2",null,[a(s,{title:c("zone-ingresses.routes.item.navigation.zone-ingress-xds-config-view")},null,8,["title"])])]),default:t(()=>[_(),a(r,null,{default:t(()=>[a(d,{resource:"Zone",src:`/zone-ingresses/${e.params.zoneIngress}/data-path/xds`,query:e.params.codeSearch,"is-filter-mode":e.params.codeFilter==="true","is-reg-exp-mode":e.params.codeRegExp==="true",onQueryChange:o=>e.update({codeSearch:o}),onFilterModeChange:o=>e.update({codeFilter:o}),onRegExpModeChange:o=>e.update({codeRegExp:o})},null,8,["src","query","is-filter-mode","is-reg-exp-mode","onQueryChange","onFilterModeChange","onRegExpModeChange"])]),_:2},1024)]),_:2},1024)]),_:1})}}});export{B as default};