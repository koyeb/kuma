"use strict";(self["webpackChunkkuma_gui"]=self["webpackChunkkuma_gui"]||[]).push([[137],{44006:function(e,t,a){a.d(t,{Z:function(){return c}});var n=function(){var e=this,t=e._self._c;return t("div",{staticClass:"wizard-switcher"},[t("KEmptyState",{ref:"emptyState",staticClass:"my-6 wizard-empty-state",attrs:{"cta-is-hidden":"","is-error":!e.environment},scopedSlots:e._u(["kubernetes"===e.environment||"universal"===e.environment?{key:"title",fn:function(){return[e._v(" Running on "),t("span",{staticClass:"capitalize"},[e._v(e._s(e.environment))])]},proxy:!0}:null,{key:"message",fn:function(){return["kubernetes"===e.environment?t("div",[e.$route.name===e.wizardRoutes.kubernetes?t("div",[t("p",[e._v(" We have detected that you are running on a "),t("strong",[e._v("Kubernetes environment")]),e._v(", and we are going to be showing you instructions for Kubernetes unless you decide to visualize the instructions for Universal. ")]),t("p",[t("KButton",{attrs:{to:{name:e.wizardRoutes.universal},appearance:"secondary"}},[e._v(" Switch to Universal instructions ")])],1)]):e.$route.name===e.wizardRoutes.universal?t("div",[t("p",[e._v(" We have detected that you are running on a "),t("strong",[e._v("Kubernetes environment")]),e._v(", but you are viewing instructions for Universal. ")]),t("p",[t("KButton",{attrs:{to:{name:e.wizardRoutes.kubernetes},appearance:"secondary"}},[e._v(" Switch back to Kubernetes instructions ")])],1)]):e._e()]):"universal"===e.environment?t("div",[e.$route.name===e.wizardRoutes.kubernetes?t("div",[t("p",[e._v(" We have detected that you are running on a "),t("strong",[e._v("Universal environment")]),e._v(", but you are viewing instructions for Kubernetes. ")]),t("p",[t("KButton",{attrs:{to:{name:e.wizardRoutes.universal},appearance:"secondary"}},[e._v(" Switch back to Universal instructions ")])],1)]):e.$route.name===e.wizardRoutes.universal?t("div",[t("p",[e._v(" We have detected that you are running on a "),t("strong",[e._v("Universal environment")]),e._v(", and we are going to be showing you instructions for Universal unless you decide to visualize the instructions for Kubernetes. ")]),t("p",[t("KButton",{attrs:{to:{name:e.wizardRoutes.kubernetes},appearance:"secondary"}},[e._v(" Switch to Kubernetes instructions ")])],1)]):e._e()]):e._e()]},proxy:!0}],null,!0)})],1)},s=[],r=a(20629),i={name:"EnvironmentSwitcher",data(){return{wizardRoutes:{kubernetes:"kubernetes-dataplane",universal:"universal-dataplane"}}},computed:{...(0,r.Se)({environment:"config/getEnvironment"}),instructionsCtaText(){return"universal"===this.environment?"Switch to Kubernetes instructions":"Switch to Universal instructions"},instructionsCtaRoute(){return"kubernetes"===this.environment?{name:"universal-dataplane"}:{name:"kubernetes-dataplane"}}}},o=i,l=a(1001),d=(0,l.Z)(o,n,s,!1,null,null,null),c=d.exports},80765:function(e,t,a){a.r(t),a.d(t,{default:function(){return w}});var n=function(){var e=this,t=e._self._c;return t("div",{staticClass:"wizard"},[t("div",{staticClass:"wizard__content"},[t("StepSkeleton",{attrs:{steps:e.steps,"sidebar-content":e.sidebarContent,"footer-enabled":!1===e.hideScannerSiblings,"next-disabled":e.nextDisabled},scopedSlots:e._u([{key:"general",fn:function(){return[t("h3",[e._v(" Create Kubernetes Dataplane ")]),t("p",[e._v(" Welcome to the wizard to create a new Dataplane resource in "+e._s(e.title)+". We will be providing you with a few steps that will get you started. ")]),t("p",[e._v(" As you know, the "+e._s(e.productName)+" GUI is read-only. ")]),t("h3",[e._v(" To get started, please select on what Mesh you would like to add the Dataplane: ")]),t("p",[e._v(" If you've got an existing Mesh that you would like to associate with your Dataplane, you can select it below, or create a new one using our Mesh Wizard. ")]),t("small",[e._v("Would you like to see instructions for Universal? Use sidebar to change wizard!")]),t("KCard",{staticClass:"my-6",attrs:{"has-shadow":""},scopedSlots:e._u([{key:"body",fn:function(){return[t("FormFragment",{attrs:{title:"Choose a Mesh","for-attr":"dp-mesh","all-inline":""}},[t("div",[t("select",{directives:[{name:"model",rawName:"v-model",value:e.validate.meshName,expression:"validate.meshName"}],staticClass:"k-input w-100",attrs:{id:"dp-mesh"},on:{change:function(t){var a=Array.prototype.filter.call(t.target.options,(function(e){return e.selected})).map((function(e){var t="_value"in e?e._value:e.value;return t}));e.$set(e.validate,"meshName",t.target.multiple?a:a[0])}}},[t("option",{attrs:{disabled:"",value:""}},[e._v(" Select an existing Mesh… ")]),e._l(e.meshes.items,(function(a){return t("option",{key:a.name,domProps:{value:a.name}},[e._v(" "+e._s(a.name)+" ")])}))],2)]),t("div",[t("label",{staticClass:"k-input-label mr-4"},[e._v(" or ")]),t("KButton",{attrs:{to:{name:"create-mesh"},appearance:"secondary"}},[e._v(" Create a new Mesh ")])],1)])]},proxy:!0}])})]},proxy:!0},{key:"scope-settings",fn:function(){return[t("h3",[e._v(" Setup Dataplane Mode ")]),t("p",[e._v(" You can create a data plane for a service or a data plane for a Gateway. ")]),t("KCard",{staticClass:"my-6",attrs:{"has-shadow":""},scopedSlots:e._u([{key:"body",fn:function(){return[t("FormFragment",{attrs:{"all-inline":"","equal-cols":"","hide-label-col":""}},[t("label",{attrs:{for:"service-dataplane"}},[t("input",{directives:[{name:"model",rawName:"v-model",value:e.validate.k8sDataplaneType,expression:"validate.k8sDataplaneType"}],staticClass:"k-input",attrs:{id:"service-dataplane",type:"radio",name:"dataplane-type",value:"dataplane-type-service",checked:""},domProps:{checked:e._q(e.validate.k8sDataplaneType,"dataplane-type-service")},on:{change:function(t){return e.$set(e.validate,"k8sDataplaneType","dataplane-type-service")}}}),t("span",[e._v(" Service Dataplane ")])]),t("label",{attrs:{for:"ingress-dataplane"}},[t("input",{directives:[{name:"model",rawName:"v-model",value:e.validate.k8sDataplaneType,expression:"validate.k8sDataplaneType"}],staticClass:"k-input",attrs:{id:"ingress-dataplane",type:"radio",name:"dataplane-type",value:"dataplane-type-ingress",disabled:""},domProps:{checked:e._q(e.validate.k8sDataplaneType,"dataplane-type-ingress")},on:{change:function(t){return e.$set(e.validate,"k8sDataplaneType","dataplane-type-ingress")}}}),t("span",[e._v(" Ingress Dataplane ")])])])]},proxy:!0}])}),"dataplane-type-service"===e.validate.k8sDataplaneType?t("div",[t("p",[e._v(" Should the data plane be added for an entire Namespace and all of its services, or for specific individual services in any namespace? ")]),t("KCard",{staticClass:"my-6",attrs:{"has-shadow":""},scopedSlots:e._u([{key:"body",fn:function(){return[t("FormFragment",{attrs:{"all-inline":"","equal-cols":"","hide-label-col":""}},[t("label",{attrs:{for:"k8s-services-all"}},[t("input",{directives:[{name:"model",rawName:"v-model",value:e.validate.k8sServices,expression:"validate.k8sServices"}],staticClass:"k-input",attrs:{id:"k8s-services-all",type:"radio",name:"k8s-services",value:"all-services",checked:""},domProps:{checked:e._q(e.validate.k8sServices,"all-services")},on:{change:function(t){return e.$set(e.validate,"k8sServices","all-services")}}}),t("span",[e._v(" All Services in Namespace ")])]),t("label",{attrs:{for:"k8s-services-individual"}},[t("input",{directives:[{name:"model",rawName:"v-model",value:e.validate.k8sServices,expression:"validate.k8sServices"}],staticClass:"k-input",attrs:{id:"k8s-services-individual",type:"radio",name:"k8s-services",value:"individual-services",disabled:""},domProps:{checked:e._q(e.validate.k8sServices,"individual-services")},on:{change:function(t){return e.$set(e.validate,"k8sServices","individual-services")}}}),t("span",[e._v(" Individual Services ")])])])]},proxy:!0}],null,!1,2127996134)}),"individual-services"===e.validate.k8sServices?t("KCard",{staticClass:"my-6",attrs:{"has-shadow":""},scopedSlots:e._u([{key:"body",fn:function(){return[t("FormFragment",{attrs:{title:"Deployments","for-attr":"k8s-deployment-selection"}},[t("input",{directives:[{name:"model",rawName:"v-model",value:e.validate.k8sServiceDeploymentSelection,expression:"validate.k8sServiceDeploymentSelection"}],staticClass:"k-input w-100",attrs:{id:"k8s-service-deployment-new",type:"text",placeholder:"your-new-deployment",required:""},domProps:{value:e.validate.k8sServiceDeploymentSelection},on:{input:function(t){t.target.composing||e.$set(e.validate,"k8sServiceDeploymentSelection",t.target.value)}}})])]},proxy:!0}],null,!1,1626108368)}):e._e(),t("KCard",{staticClass:"my-6",attrs:{"has-shadow":""},scopedSlots:e._u([{key:"body",fn:function(){return[t("FormFragment",{attrs:{title:"Namespace","for-attr":"k8s-namespace-selection"}},[t("input",{directives:[{name:"model",rawName:"v-model",value:e.validate.k8sNamespaceSelection,expression:"validate.k8sNamespaceSelection"}],staticClass:"k-input w-100",attrs:{id:"k8s-namespace-new",type:"text",placeholder:"your-namespace",required:""},domProps:{value:e.validate.k8sNamespaceSelection},on:{input:function(t){t.target.composing||e.$set(e.validate,"k8sNamespaceSelection",t.target.value)}}})])]},proxy:!0}],null,!1,771225282)})],1):e._e(),"dataplane-type-ingress"===e.validate.k8sDataplaneType?t("div",[t("p",[e._v(" "+e._s(e.title)+" natively supports the Kong Ingress. Do you want to deploy Kong or another Ingress? ")]),t("KCard",{staticClass:"my-6",attrs:{"has-shadow":""},scopedSlots:e._u([{key:"body",fn:function(){return[t("FormFragment",{attrs:{"all-inline":"","equal-cols":"","hide-label-col":""}},[t("label",{attrs:{for:"k8s-ingress-kong"}},[t("input",{directives:[{name:"model",rawName:"v-model",value:e.validate.k8sIngressBrand,expression:"validate.k8sIngressBrand"}],staticClass:"k-input",attrs:{id:"k8s-ingress-kong",type:"radio",name:"k8s-ingress-brand",value:"kong-ingress",checked:""},domProps:{checked:e._q(e.validate.k8sIngressBrand,"kong-ingress")},on:{change:function(t){return e.$set(e.validate,"k8sIngressBrand","kong-ingress")}}}),t("span",[e._v(" Kong Ingress ")])]),t("label",{attrs:{for:"k8s-ingress-other"}},[t("input",{directives:[{name:"model",rawName:"v-model",value:e.validate.k8sIngressBrand,expression:"validate.k8sIngressBrand"}],staticClass:"k-input",attrs:{id:"k8s-ingress-other",type:"radio",name:"k8s-ingress-brand",value:"other-ingress"},domProps:{checked:e._q(e.validate.k8sIngressBrand,"other-ingress")},on:{change:function(t){return e.$set(e.validate,"k8sIngressBrand","other-ingress")}}}),t("span",[e._v(" Other Ingress ")])])])]},proxy:!0}],null,!1,1060751940)}),t("KCard",{staticClass:"my-6",attrs:{"has-shadow":""},scopedSlots:e._u([{key:"body",fn:function(){return[t("FormFragment",{attrs:{title:"Deployments","for-attr":"k8s-deployment-selection"}},[t("input",{directives:[{name:"model",rawName:"v-model",value:e.validate.k8sIngressDeployment,expression:"validate.k8sIngressDeployment"}],staticClass:"k-input w-100",attrs:{id:"k8s-ingress-deployment-new",type:"text",placeholder:"your-deployment",required:""},domProps:{value:e.validate.k8sIngressDeployment},on:{input:function(t){t.target.composing||e.$set(e.validate,"k8sIngressDeployment",t.target.value)}}})])]},proxy:!0}],null,!1,1817964619)}),"other-ingress"===e.validate.k8sIngressBrand?t("KAlert",{attrs:{appearance:"info"},scopedSlots:e._u([{key:"alertMessage",fn:function(){return[t("p",[e._v(' Please go ahead and deploy the Ingress first, then restart this wizard and select "Existing Ingress". ')])]},proxy:!0}],null,!1,1402213972)}):e._e()],1):e._e()]},proxy:!0},{key:"complete",fn:function(){return[e.validate.meshName?t("div",[!1===e.hideScannerSiblings?t("div",[t("h3",[e._v(" Auto-Inject DPP ")]),t("p",[e._v(" You can now execute the following commands to automatically inject the sidecar proxy in every Pod, and by doing so creating the Dataplane. ")]),t("TabsWidget",{attrs:{loaders:!1,tabs:e.tabs,"initial-tab-override":"kubernetes"},scopedSlots:e._u([{key:"kubernetes",fn:function(){return[t("CodeView",{attrs:{title:"Kubernetes","copy-button-text":"Copy Command to Clipboard",lang:"bash",content:e.codeOutput}})]},proxy:!0}],null,!1,525752398)})],1):e._e(),t("EntityScanner",{attrs:{"loader-function":e.scanForEntity,"should-start":!0,"has-error":e.scanError,"can-complete":e.scanFound},on:{hideSiblings:e.hideSiblings},scopedSlots:e._u([{key:"loading-title",fn:function(){return[t("h3",[e._v("Searching…")])]},proxy:!0},{key:"loading-content",fn:function(){return[t("p",[e._v("We are looking for your dataplane.")])]},proxy:!0},{key:"complete-title",fn:function(){return[t("h3",[e._v("Done!")])]},proxy:!0},{key:"complete-content",fn:function(){return[t("p",[e._v(" Your Dataplane "),e.validate.k8sNamespaceSelection?t("strong",[e._v(" "+e._s(e.validate.k8sNamespaceSelection)+" ")]):e._e(),e._v(" was found! ")]),t("p",[e._v(" Proceed to the next step where we will show you your new Dataplane. ")]),t("p",[t("KButton",{attrs:{appearance:"primary"},on:{click:e.compeleteDataPlaneSetup}},[e._v(" View Your Dataplane ")])],1)]},proxy:!0},{key:"error-title",fn:function(){return[t("h3",[e._v("Mesh not found")])]},proxy:!0},{key:"error-content",fn:function(){return[t("p",[e._v("We were unable to find your mesh.")])]},proxy:!0}],null,!1,2302604054)})],1):t("KAlert",{attrs:{appearance:"danger"},scopedSlots:e._u([{key:"alertMessage",fn:function(){return[t("p",[e._v(" Please return to the first step and make sure to select an existing Mesh, or create a new one. ")])]},proxy:!0}])})]},proxy:!0},{key:"dataplane",fn:function(){return[t("h3",[e._v("Dataplane")]),t("p",[e._v(" In "+e._s(e.title)+", a Dataplane resource represents a data plane proxy running alongside one of your services. Data plane proxies can be added in any Mesh that you may have created, and in Kubernetes, they will be auto-injected by "+e._s(e.title)+". ")])]},proxy:!0},{key:"example",fn:function(){return[t("h3",[e._v("Example")]),t("p",[e._v(" Below is an example of a Dataplane resource output: ")]),t("code",{staticClass:"block"},[t("pre",[e._v("apiVersion: 'kuma.io/v1alpha1'\nkind: Dataplane\nmesh: default\nmetadata:\n  name: dp-echo-1\n  annotations:\n    kuma.io/sidecar-injection: enabled\n    kuma.io/mesh: default\nnetworking:\n  address: 10.0.0.1\n  inbound:\n  - port: 10000\n    servicePort: 9000\n    tags:\n      kuma.io/service: echo")])])]},proxy:!0},{key:"switch",fn:function(){return[t("EnvironmentSwitcher")]},proxy:!0}])})],1)])},s=[],r=a(20629),i=a(17463),o=a(53419),l=a(88523),d=a(5035),c=a(7001),u=a(5872),p=a(44006),v=a(22330),m=a(69328),h=a(89938),y=a.n(h),k=a(45689),g={name:"DataplaneWizardKubernetes",metaInfo:{title:"Create a new Dataplane on Kubernetes"},components:{FormFragment:d.Z,TabsWidget:c.Z,StepSkeleton:u.Z,EnvironmentSwitcher:p.Z,CodeView:v.Z,EntityScanner:m.Z},mixins:[l.Z],data(){return{productName:k.sG,schema:y(),steps:[{label:"General",slug:"general"},{label:"Scope Settings",slug:"scope-settings"},{label:"Install",slug:"complete"}],tabs:[{hash:"#kubernetes",title:"Kubernetes"}],sidebarContent:[{name:"dataplane"},{name:"example"},{name:"switch"}],startScanner:!1,scanFound:!1,hideScannerSiblings:!1,scanError:!1,isComplete:!1,validate:{meshName:"",k8sDataplaneType:"dataplane-type-service",k8sServices:"all-services",k8sNamespace:"",k8sNamespaceSelection:"",k8sServiceDeployment:"",k8sServiceDeploymentSelection:"",k8sIngressDeployment:"",k8sIngressDeploymentSelection:"",k8sIngressType:"",k8sIngressBrand:"kong-ingress",k8sIngressSelection:""}}},computed:{...(0,r.Se)({title:"config/getTagline",version:"config/getVersion",environment:"config/getEnvironment",meshes:"getMeshList"}),dataplaneUrl(){const e=this.validate;return!(!e.meshName||!e.k8sNamespaceSelection)&&{name:"dataplanes",params:{mesh:e.meshName}}},codeOutput(){const e=Object.assign({},this.schema),t=this.validate.k8sNamespaceSelection;if(!t)return;e.metadata.name=t,e.metadata.namespace=t,e.metadata.annotations["kuma.io/mesh"]=this.validate.meshName;const a=`" | kubectl apply -f - && kubectl delete pod --all -n ${t}`,n=this.formatForCLI(e,a);return n},nextDisabled(){const{k8sNamespaceSelection:e,meshName:t}=this.validate;return!t.length||"1"===this.$route.query.step&&!e}},watch:{"validate.k8sNamespaceSelection"(e){this.validate.k8sNamespaceSelection=(0,o.GL)(e)},$route(){const e=this.$route.query.step;1===e&&(this.validate.k8sNamespaceSelection?this.nextDisabled=!1:this.nextDisabled=!0)}},methods:{hideSiblings(){this.hideScannerSiblings=!0},scanForEntity(){const e=this.validate,t=e.meshName,a=this.validate.k8sNamespaceSelection;this.scanComplete=!1,this.scanError=!1,t&&a&&i.Z.getDataplaneFromMesh({mesh:t,name:a}).then((e=>{e&&e.name.length>0?(this.isRunning=!0,this.scanFound=!0):this.scanError=!0})).catch((e=>{this.scanError=!0,console.error(e)})).finally((()=>{this.scanComplete=!0}))},compeleteDataPlaneSetup(){this.$store.dispatch("updateSelectedMesh",this.validate.meshName),this.$router.push({name:"dataplanes",params:{mesh:this.validate.meshName}})}}},f=g,_=a(1001),b=(0,_.Z)(f,n,s,!1,null,"7c87db2a",null),w=b.exports},88523:function(e,t,a){var n=a(73570),s=a.n(n);t["Z"]={methods:{formatForCLI(e,t='" | kumactl apply -f -'){const a='echo "',n=s()(e);return`${a}${n}${t}`}}}},89938:function(e){e.exports={apiVersion:"v1",kind:"Namespace",metadata:{name:null,namespace:null,annotations:{"kuma.io/sidecar-injection":"enabled","kuma.io/mesh":null}}}}}]);