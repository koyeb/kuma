import{d as A,L as C,q as $,c as B,o as n,e as l,k as h,n as x,g as i,h as b,w as d,f as p,T as E,y as j,M as q,a as f,t as u,b as c,F as T,j as S,p as N,m as V}from"./index-bec4d193.js";import{g as O,i as H,r as D,A as L}from"./RouteView.vue_vue_type_script_setup_true_lang-19528cd9.js";import{l as M,Z as P}from"./kongponents.es-b2b550a8.js";import{a as y,D as w}from"./DefinitionListItem-d3bb1a9c.js";const F=["aria-expanded"],R={key:0,class:"accordion-item-content","data-testid":"accordion-item-content"},Z=A({__name:"AccordionItem",setup(t){const e=C("parentAccordion"),a=$(null),r=B(()=>e===void 0?!1:e.multipleOpen&&Array.isArray(e.active.value)&&a.value!==null?e.active.value.includes(a.value):a.value===e.active.value);e!==void 0&&(a.value=e.count.value++);function _(){r.value?I():o()}function I(){e!==void 0&&(e.multipleOpen&&Array.isArray(e.active.value)&&a.value!==null?e.active.value.splice(e.active.value.indexOf(a.value),1):e.active.value=null)}function o(){e!==void 0&&(e.multipleOpen&&Array.isArray(e.active.value)&&a.value!==null?e.active.value.push(a.value):e.active.value=a.value)}function v(s){s instanceof HTMLElement&&(s.style.height=`${s.scrollHeight}px`)}function m(s){s instanceof HTMLElement&&(s.style.height="auto")}return(s,k)=>(n(),l("li",{class:j(["accordion-item",{active:r.value}])},[h("button",{class:"accordion-item-header",type:"button","aria-expanded":r.value?"true":"false","data-testid":"accordion-item-button",onClick:_},[x(s.$slots,"accordion-header",{},void 0,!0)],8,F),i(),b(E,{name:"accordion",onEnter:v,onAfterEnter:m,onBeforeLeave:v},{default:d(()=>[r.value?(n(),l("div",R,[x(s.$slots,"accordion-content",{},void 0,!0)])):p("",!0)]),_:3})],2))}});const ce=O(Z,[["__scopeId","data-v-dfd99690"]]),z={class:"accordion-list"},G=A({__name:"AccordionList",props:{initiallyOpen:{type:[Number,Array],required:!1,default:null},multipleOpen:{type:Boolean,required:!1,default:!1}},setup(t){const e=t,a=$(0),r=$(e.initiallyOpen!==null?e.initiallyOpen:e.multipleOpen?[]:null);return q("parentAccordion",{multipleOpen:e.multipleOpen,active:r,count:a}),(_,I)=>(n(),l("ul",z,[x(_.$slots,"default",{},void 0,!0)]))}});const le=O(G,[["__scopeId","data-v-53d92d22"]]),U=t=>(N("data-v-321555ca"),t=t(),V(),t),J={key:0},K=U(()=>h("h5",{class:"overview-tertiary-title"},`
        General Information:
      `,-1)),Q={key:1,class:"columns mt-4",style:{"--columns":"4"}},W={key:0},X={class:"overview-tertiary-title"},Y=A({__name:"SubscriptionDetails",props:{details:{type:Object,required:!0},isDiscoverySubscription:{type:Boolean,default:!1}},setup(t){const e=t,{t:a}=H(),r=B(()=>{var v,m;let o;if(e.isDiscoverySubscription){const{lastUpdateTime:s,total:k,...g}=e.details.status;o=g}return(v=e.details.status)!=null&&v.stat&&(o=(m=e.details.status)==null?void 0:m.stat),o});function _(o){return o?parseInt(o,10).toLocaleString("en").toString():"0"}function I(o){return o==="--"?"error calculating":o}return(o,v)=>(n(),l("div",null,[t.details.globalInstanceId||t.details.controlPlaneInstanceId||t.details.connectTime||t.details.disconnectTime?(n(),l("div",J,[K,i(),b(w,null,{default:d(()=>[t.details.globalInstanceId?(n(),f(y,{key:0,term:c(a)("http.api.property.globalInstanceId")},{default:d(()=>[i(u(t.details.globalInstanceId),1)]),_:1},8,["term"])):p("",!0),i(),t.details.controlPlaneInstanceId?(n(),f(y,{key:1,term:c(a)("http.api.property.controlPlaneInstanceId")},{default:d(()=>[i(u(t.details.controlPlaneInstanceId),1)]),_:1},8,["term"])):p("",!0),i(),t.details.connectTime?(n(),f(y,{key:2,term:c(a)("http.api.property.connectTime")},{default:d(()=>[i(u(c(D)(t.details.connectTime)),1)]),_:1},8,["term"])):p("",!0),i(),t.details.disconnectTime?(n(),f(y,{key:3,term:c(a)("http.api.property.disconnectTime")},{default:d(()=>[i(u(c(D)(t.details.disconnectTime)),1)]),_:1},8,["term"])):p("",!0)]),_:1})])):p("",!0),i(),r.value?(n(),l("div",Q,[(n(!0),l(T,null,S(r.value,(m,s)=>(n(),l(T,{key:s},[Object.keys(m).length>0?(n(),l("div",W,[h("h6",X,u(c(a)(`http.api.property.${s}`))+`:
          `,1),i(),b(w,null,{default:d(()=>[(n(!0),l(T,null,S(m,(k,g)=>(n(),f(y,{key:g,term:c(a)(`http.api.property.${g}`)},{default:d(()=>[i(u(I(_(k))),1)]),_:2},1032,["term"]))),128))]),_:2},1024)])):p("",!0)],64))),128))])):(n(),f(c(P),{key:2,appearance:"info",class:"mt-4"},{alertIcon:d(()=>[b(c(M),{icon:"portal"})]),alertMessage:d(()=>[i(`
        There are no subscription statistics for `),h("strong",null,u(t.details.id),1)]),_:1}))]))}});const re=O(Y,[["__scopeId","data-v-321555ca"]]),ee={class:"text-lg font-medium"},te={class:"color-green-500"},ae={key:0,class:"ml-4 color-red-600"},de=A({__name:"SubscriptionHeader",props:{details:{type:Object,required:!0}},setup(t){const e=t;return(a,r)=>(n(),l("h4",ee,[h("span",te,`
      Connect time: `+u(c(L)(e.details.connectTime)),1),i(),e.details.disconnectTime?(n(),l("span",ae,`
      Disconnect time: `+u(c(L)(e.details.disconnectTime)),1)):p("",!0)]))}});export{ce as A,re as S,de as _,le as a};