import{E as V}from"./ErrorBlock-567378ca.js";import{_ as g}from"./LoadingBlock.vue_vue_type_script_setup_true_lang-52bdda9b.js";import{_ as k}from"./ResourceCodeBlock.vue_vue_type_style_index_0_lang-60c9f4a9.js";import{d as v,a as o,o as t,b as r,w as e,e as s,m as x,f as B,p as E}from"./index-5d5446a4.js";import"./index-fce48c05.js";import"./TextWithCopyButton-1669005d.js";import"./CopyButton-b62a1694.js";import"./WarningIcon.vue_vue_type_script_setup_true_lang-e798630e.js";import"./CodeBlock.vue_vue_type_style_index_0_lang-be5bb8ae.js";import"./uniqueId-90cc9b93.js";import"./toYaml-4e00099e.js";const F=v({__name:"MeshConfigView",setup(N){return(R,$)=>{const l=o("RouteTitle"),a=o("DataSource"),u=o("KCard"),f=o("AppView"),d=o("RouteView");return t(),r(d,{name:"mesh-config-view","data-testid":"mesh-config-view",params:{mesh:""}},{default:e(({route:m,t:h})=>[s(f,null,{title:e(()=>[x("h2",null,[s(l,{title:h("meshes.routes.item.navigation.mesh-config-view")},null,8,["title"])])]),default:e(()=>[B(),s(u,null,{default:e(()=>[s(a,{src:`/meshes/${m.params.mesh}`},{default:e(({data:i,error:c})=>[c!==void 0?(t(),r(V,{key:0,error:c},null,8,["error"])):i===void 0?(t(),r(g,{key:1})):(t(),r(k,{key:2,resource:i.config},{default:e(({copy:p,copying:w})=>[w?(t(),r(a,{key:0,src:`/meshes/${m.params.mesh}/as/kubernetes?no-store`,onChange:n=>{p(_=>_(n))},onError:n=>{p((_,C)=>C(n))}},null,8,["src","onChange","onError"])):E("",!0)]),_:2},1032,["resource"]))]),_:2},1032,["src"])]),_:2},1024)]),_:2},1024)]),_:1})}}});export{F as default};