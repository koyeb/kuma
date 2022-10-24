import{C as y,cs as T,ct as p,o as s,k as n,l as a,b as r,t as c,y as i,F as h,n as b,c as w,w as I,a as x,j as f,A as C,B as R,cu as B}from"./index.ea0d4a24.js";const V={name:"SubscriptionDetails",props:{details:{type:Object,required:!0},isDiscoverySubscription:{type:Boolean,default:!1}},computed:{detailsIterator(){var e;if(this.isDiscoverySubscription){const{lastUpdateTime:_,total:t,...d}=this.details.status;return d}return(e=this.details.status)==null?void 0:e.stat}},methods:{formatValue(e){return e?parseInt(e,10).toLocaleString("en").toString():0},readableDate(e){return T(e)},humanReadable(e){return p(e)},formatError(e){return e==="--"?"error calculating":e}}},l=e=>(C("data-v-2dfb031a"),e=e(),R(),e),F={key:0},K=l(()=>a("h5",{class:"overview-tertiary-title"}," General Information: ",-1)),L={key:0},N=l(()=>a("strong",null,"Global Instance ID:",-1)),j={class:"mono"},A={key:1},E=l(()=>a("strong",null,"Control Plane Instance ID:",-1)),P={class:"mono"},q={key:2},G=l(()=>a("strong",null,"Last Connected:",-1)),H={key:3},O=l(()=>a("strong",null,"Last Disconnected:",-1)),M={key:1},U={class:"overview-stat-grid"},W={class:"overview-tertiary-title"},z={class:"mono"};function J(e,_,t,d,g,o){const D=f("KIcon"),k=f("KAlert");return s(),n("div",null,[t.details.globalInstanceId||t.details.connectTime||t.details.disconnectTime?(s(),n("div",F,[K,a("ul",null,[t.details.globalInstanceId?(s(),n("li",L,[N,r("\xA0 "),a("span",j,c(t.details.globalInstanceId),1)])):i("",!0),t.details.controlPlaneInstanceId?(s(),n("li",A,[E,r("\xA0 "),a("span",P,c(t.details.controlPlaneInstanceId),1)])):i("",!0),t.details.connectTime?(s(),n("li",q,[G,r("\xA0 "+c(o.readableDate(t.details.connectTime)),1)])):i("",!0),t.details.disconnectTime?(s(),n("li",H,[O,r("\xA0 "+c(o.readableDate(t.details.disconnectTime)),1)])):i("",!0)])])):i("",!0),o.detailsIterator?(s(),n("div",M,[a("ul",U,[(s(!0),n(h,null,b(o.detailsIterator,(S,u)=>(s(),n("li",{key:u},[a("h6",W,c(o.humanReadable(u))+": ",1),a("ul",null,[(s(!0),n(h,null,b(S,(v,m)=>(s(),n("li",{key:m},[a("strong",null,c(o.humanReadable(m))+":",1),r("\xA0 "),a("span",z,c(o.formatError(o.formatValue(v))),1)]))),128))])]))),128))])])):(s(),w(k,{key:2,appearance:"info",class:"mt-4"},{alertIcon:I(()=>[x(D,{icon:"portal"})]),alertMessage:I(()=>[r(" There are no subscription statistics for "),a("strong",null,c(t.details.id),1)]),_:1}))])}const te=y(V,[["render",J],["__scopeId","data-v-2dfb031a"]]),Q={name:"SubscriptionHeader",props:{details:{type:Object,required:!0}},methods:{rawReadableDateFilter(e){return B(e)}}},X={class:"text-lg font-medium"},Y={class:"color-green-400"},Z={key:0,class:"ml-4 color-red-400"};function $(e,_,t,d,g,o){return s(),n("h4",X,[a("span",Y," Connect time: "+c(o.rawReadableDateFilter(t.details.connectTime)),1),t.details.disconnectTime?(s(),n("span",Z," Disconnect time: "+c(o.rawReadableDateFilter(t.details.disconnectTime)),1)):i("",!0)])}const ae=y(Q,[["render",$]]);export{te as S,ae as a};
