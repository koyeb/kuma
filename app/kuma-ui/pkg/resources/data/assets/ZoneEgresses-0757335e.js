import{d as $,l as M,m as Q,r as o,n as R,p as U,s as G,t as H,o as i,c as v,g as d,a as u,w as s,k as _,e as E,u as N,E as K,x as b,P as q,y as z,F as B,z as C,B as W,C as X}from"./index-3d59543a.js";import{D as Y}from"./DataOverview-2890198a.js";import{E as k}from"./EnvoyData-2003e9c8.js";import{_ as j}from"./LabelList.vue_vue_type_style_index_0_lang-38bc105d.js";import{_ as J,S as ee}from"./SubscriptionHeader.vue_vue_type_script_setup_true_lang-86ad7222.js";import{T as ae}from"./TabsWidget-29ed07c3.js";import{Q as I}from"./QueryParameter-70743f73.js";import"./EmptyBlock.vue_vue_type_script_setup_true_lang-f6291874.js";import"./ErrorBlock-efee1ec6.js";import"./LoadingBlock.vue_vue_type_script_setup_true_lang-08383420.js";import"./TagList-9abb7297.js";import"./StatusBadge-091066a4.js";import"./CodeBlock.vue_vue_type_style_index_0_lang-3e8641de.js";const se={class:"zoneegresses"},te={class:"kcard-stack"},ne={class:"kcard-border"},re={class:"kcard-border"},oe={class:"entity-heading"},le={key:0},be=$({__name:"ZoneEgresses",props:{selectedZoneEgressName:{type:String,required:!1,default:null},offset:{type:Number,required:!1,default:0}},setup(O){const p=O,w=M(),L={title:"No Data",message:"There are no Zone Egresses present."},V=[{hash:"#overview",title:"Overview"},{hash:"#insights",title:"Zone Egress Insights"},{hash:"#xds-configuration",title:"XDS Configuration"},{hash:"#envoy-stats",title:"Stats"},{hash:"#envoy-clusters",title:"Clusters"}],f=Q(),m=o(!0),c=o(!1),g=o(null),y=o({headers:[{label:"Status",key:"status"},{label:"Name",key:"name"}],data:[]}),l=o(null),x=o([]),D=o(null),S=o([]),A=o(p.offset);R(()=>f.params.mesh,function(){f.name==="zoneegresses"&&(m.value=!0,c.value=!1,g.value=null,h(0))}),U(function(){h(p.offset)});async function h(a){A.value=a,I.set("offset",a>0?a:null),m.value=!0,c.value=!1;const t=f.query.ns||null,r=q;try{const{data:e,next:n}=await F(t,r,a);D.value=n,e.length?(c.value=!1,x.value=e,T({name:p.selectedZoneEgressName??e[0].name}),y.value.data=e.map(Z=>{const P=G(Z.zoneEgressInsight??{});return{...Z,status:P}})):(y.value.data=[],c.value=!0)}catch(e){e instanceof Error?g.value=e:console.error(e),c.value=!0}finally{m.value=!1}}function T({name:a}){var e;const t=x.value.find(n=>n.name===a),r=((e=t==null?void 0:t.zoneEgressInsight)==null?void 0:e.subscriptions)??[];S.value=Array.from(r).reverse(),l.value=H(t,["type","name"]),I.set("zoneEgress",a)}async function F(a,t,r){if(a)return{data:[await w.getZoneEgressOverview({name:a},{size:t,offset:r})],next:null};{const{items:e,next:n}=await w.getAllZoneEgressOverviews({size:t,offset:r});return{data:e??[],next:n}}}return(a,t)=>{var r;return i(),v("div",se,[d("div",te,[d("div",ne,[u(Y,{"selected-entity-name":(r=l.value)==null?void 0:r.name,"page-size":N(q),"is-loading":m.value,error:g.value,"empty-state":L,"table-data":y.value,"table-data-is-empty":c.value,next:D.value,"page-offset":A.value,onTableAction:T,onLoadData:h},{additionalControls:s(()=>[a.$route.query.ns?(i(),_(N(K),{key:0,class:"back-button",appearance:"primary",icon:"arrowLeft",to:{name:"zoneegresses"}},{default:s(()=>[E(`
              View all
            `)]),_:1})):b("",!0)]),_:1},8,["selected-entity-name","page-size","is-loading","error","table-data","table-data-is-empty","next","page-offset"])]),E(),d("div",re,[c.value===!1&&l.value!==null?(i(),_(ae,{key:0,"has-error":g.value!==null,"is-loading":m.value,tabs:V},{tabHeader:s(()=>[d("h1",oe,`
              Zone Egress: `+z(l.value.name),1)]),overview:s(()=>[u(j,null,{default:s(()=>[d("div",null,[d("ul",null,[(i(!0),v(B,null,C(l.value,(e,n)=>(i(),v("li",{key:n},[e?(i(),v("h4",le,z(n),1)):b("",!0),E(),d("p",null,z(e),1)]))),128))])])]),_:1})]),insights:s(()=>[u(X,{"initially-open":0},{default:s(()=>[(i(!0),v(B,null,C(S.value,(e,n)=>(i(),_(W,{key:n},{"accordion-header":s(()=>[u(J,{details:e},null,8,["details"])]),"accordion-content":s(()=>[u(ee,{details:e,"is-discovery-subscription":""},null,8,["details"])]),_:2},1024))),128))]),_:1})]),"xds-configuration":s(()=>[u(k,{"data-path":"xds","zone-egress-name":l.value.name,"query-key":"envoy-data-zone-egress"},null,8,["zone-egress-name"])]),"envoy-stats":s(()=>[u(k,{"data-path":"stats","zone-egress-name":l.value.name,"query-key":"envoy-data-zone-egress"},null,8,["zone-egress-name"])]),"envoy-clusters":s(()=>[u(k,{"data-path":"clusters","zone-egress-name":l.value.name,"query-key":"envoy-data-zone-egress"},null,8,["zone-egress-name"])]),_:1},8,["has-error","is-loading"])):b("",!0)])])])}}});export{be as default};