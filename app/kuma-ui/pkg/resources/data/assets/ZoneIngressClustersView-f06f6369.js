import{E as l}from"./EnvoyData-0ef3d659.js";import{d as m,a as e,o as u,b as _,w as o,e as t,p as d,f as g}from"./index-9fad9c9c.js";import"./index-52545d1d.js";import"./CodeBlock.vue_vue_type_style_index_0_lang-38ef20ee.js";import"./EmptyBlock.vue_vue_type_script_setup_true_lang-51396c83.js";import"./ErrorBlock-9632ec5c.js";import"./TextWithCopyButton-fdcabbe9.js";import"./CopyButton-ca05cb3e.js";import"./WarningIcon.vue_vue_type_script_setup_true_lang-6cf89112.js";import"./LoadingBlock.vue_vue_type_script_setup_true_lang-6d2f87fd.js";const S=m({__name:"ZoneIngressClustersView",setup(h){return(w,f)=>{const s=e("RouteTitle"),r=e("KCard"),a=e("AppView"),i=e("RouteView");return u(),_(i,{name:"zone-ingress-clusters-view",params:{zoneIngress:"",codeSearch:""}},{default:o(({route:n,t:c})=>[t(a,null,{title:o(()=>[d("h2",null,[t(s,{title:c("zone-ingresses.routes.item.navigation.zone-ingress-clusters-view")},null,8,["title"])])]),default:o(()=>[g(),t(r,null,{body:o(()=>[t(l,{resource:"Zone",src:`/zone-ingresses/${n.params.zoneIngress}/data-path/clusters`,query:n.params.codeSearch,onQueryChange:p=>n.update({codeSearch:p})},null,8,["src","query","onQueryChange"])]),_:2},1024)]),_:2},1024)]),_:1})}}});export{S as default};