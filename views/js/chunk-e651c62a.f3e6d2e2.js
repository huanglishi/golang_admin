(window["webpackJsonp"]=window["webpackJsonp"]||[]).push([["chunk-e651c62a"],{"00d8":function(e,t){(function(){var t="ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/",r={rotl:function(e,t){return e<<t|e>>>32-t},rotr:function(e,t){return e<<32-t|e>>>t},endian:function(e){if(e.constructor==Number)return 16711935&r.rotl(e,8)|4278255360&r.rotl(e,24);for(var t=0;t<e.length;t++)e[t]=r.endian(e[t]);return e},randomBytes:function(e){for(var t=[];e>0;e--)t.push(Math.floor(256*Math.random()));return t},bytesToWords:function(e){for(var t=[],r=0,a=0;r<e.length;r++,a+=8)t[a>>>5]|=e[r]<<24-a%32;return t},wordsToBytes:function(e){for(var t=[],r=0;r<32*e.length;r+=8)t.push(e[r>>>5]>>>24-r%32&255);return t},bytesToHex:function(e){for(var t=[],r=0;r<e.length;r++)t.push((e[r]>>>4).toString(16)),t.push((15&e[r]).toString(16));return t.join("")},hexToBytes:function(e){for(var t=[],r=0;r<e.length;r+=2)t.push(parseInt(e.substr(r,2),16));return t},bytesToBase64:function(e){for(var r=[],a=0;a<e.length;a+=3)for(var n=e[a]<<16|e[a+1]<<8|e[a+2],i=0;i<4;i++)8*a+6*i<=8*e.length?r.push(t.charAt(n>>>6*(3-i)&63)):r.push("=");return r.join("")},base64ToBytes:function(e){e=e.replace(/[^A-Z0-9+\/]/gi,"");for(var r=[],a=0,n=0;a<e.length;n=++a%4)0!=n&&r.push((t.indexOf(e.charAt(a-1))&Math.pow(2,-2*n+8)-1)<<2*n|t.indexOf(e.charAt(a))>>>6-2*n);return r}};e.exports=r})()},"044b":function(e,t){function r(e){return!!e.constructor&&"function"===typeof e.constructor.isBuffer&&e.constructor.isBuffer(e)}function a(e){return"function"===typeof e.readFloatLE&&"function"===typeof e.slice&&r(e.slice(0,0))}
/*!
 * Determine if an object is a Buffer
 *
 * @author   Feross Aboukhadijeh <https://feross.org>
 * @license  MIT
 */
e.exports=function(e){return null!=e&&(r(e)||a(e)||!!e._isBuffer)}},1608:function(e,t,r){},6821:function(e,t,r){(function(){var t=r("00d8"),a=r("9a634").utf8,n=r("044b"),i=r("9a634").bin,s=function(e,r){e.constructor==String?e=r&&"binary"===r.encoding?i.stringToBytes(e):a.stringToBytes(e):n(e)?e=Array.prototype.slice.call(e,0):Array.isArray(e)||e.constructor===Uint8Array||(e=e.toString());for(var o=t.bytesToWords(e),c=8*e.length,u=1732584193,l=-271733879,g=-1732584194,d=271733878,f=0;f<o.length;f++)o[f]=16711935&(o[f]<<8|o[f]>>>24)|4278255360&(o[f]<<24|o[f]>>>8);o[c>>>5]|=128<<c%32,o[14+(c+64>>>9<<4)]=c;var m=s._ff,p=s._gg,h=s._hh,b=s._ii;for(f=0;f<o.length;f+=16){var v=u,y=l,T=g,w=d;u=m(u,l,g,d,o[f+0],7,-680876936),d=m(d,u,l,g,o[f+1],12,-389564586),g=m(g,d,u,l,o[f+2],17,606105819),l=m(l,g,d,u,o[f+3],22,-1044525330),u=m(u,l,g,d,o[f+4],7,-176418897),d=m(d,u,l,g,o[f+5],12,1200080426),g=m(g,d,u,l,o[f+6],17,-1473231341),l=m(l,g,d,u,o[f+7],22,-45705983),u=m(u,l,g,d,o[f+8],7,1770035416),d=m(d,u,l,g,o[f+9],12,-1958414417),g=m(g,d,u,l,o[f+10],17,-42063),l=m(l,g,d,u,o[f+11],22,-1990404162),u=m(u,l,g,d,o[f+12],7,1804603682),d=m(d,u,l,g,o[f+13],12,-40341101),g=m(g,d,u,l,o[f+14],17,-1502002290),l=m(l,g,d,u,o[f+15],22,1236535329),u=p(u,l,g,d,o[f+1],5,-165796510),d=p(d,u,l,g,o[f+6],9,-1069501632),g=p(g,d,u,l,o[f+11],14,643717713),l=p(l,g,d,u,o[f+0],20,-373897302),u=p(u,l,g,d,o[f+5],5,-701558691),d=p(d,u,l,g,o[f+10],9,38016083),g=p(g,d,u,l,o[f+15],14,-660478335),l=p(l,g,d,u,o[f+4],20,-405537848),u=p(u,l,g,d,o[f+9],5,568446438),d=p(d,u,l,g,o[f+14],9,-1019803690),g=p(g,d,u,l,o[f+3],14,-187363961),l=p(l,g,d,u,o[f+8],20,1163531501),u=p(u,l,g,d,o[f+13],5,-1444681467),d=p(d,u,l,g,o[f+2],9,-51403784),g=p(g,d,u,l,o[f+7],14,1735328473),l=p(l,g,d,u,o[f+12],20,-1926607734),u=h(u,l,g,d,o[f+5],4,-378558),d=h(d,u,l,g,o[f+8],11,-2022574463),g=h(g,d,u,l,o[f+11],16,1839030562),l=h(l,g,d,u,o[f+14],23,-35309556),u=h(u,l,g,d,o[f+1],4,-1530992060),d=h(d,u,l,g,o[f+4],11,1272893353),g=h(g,d,u,l,o[f+7],16,-155497632),l=h(l,g,d,u,o[f+10],23,-1094730640),u=h(u,l,g,d,o[f+13],4,681279174),d=h(d,u,l,g,o[f+0],11,-358537222),g=h(g,d,u,l,o[f+3],16,-722521979),l=h(l,g,d,u,o[f+6],23,76029189),u=h(u,l,g,d,o[f+9],4,-640364487),d=h(d,u,l,g,o[f+12],11,-421815835),g=h(g,d,u,l,o[f+15],16,530742520),l=h(l,g,d,u,o[f+2],23,-995338651),u=b(u,l,g,d,o[f+0],6,-198630844),d=b(d,u,l,g,o[f+7],10,1126891415),g=b(g,d,u,l,o[f+14],15,-1416354905),l=b(l,g,d,u,o[f+5],21,-57434055),u=b(u,l,g,d,o[f+12],6,1700485571),d=b(d,u,l,g,o[f+3],10,-1894986606),g=b(g,d,u,l,o[f+10],15,-1051523),l=b(l,g,d,u,o[f+1],21,-2054922799),u=b(u,l,g,d,o[f+8],6,1873313359),d=b(d,u,l,g,o[f+15],10,-30611744),g=b(g,d,u,l,o[f+6],15,-1560198380),l=b(l,g,d,u,o[f+13],21,1309151649),u=b(u,l,g,d,o[f+4],6,-145523070),d=b(d,u,l,g,o[f+11],10,-1120210379),g=b(g,d,u,l,o[f+2],15,718787259),l=b(l,g,d,u,o[f+9],21,-343485551),u=u+v>>>0,l=l+y>>>0,g=g+T>>>0,d=d+w>>>0}return t.endian([u,l,g,d])};s._ff=function(e,t,r,a,n,i,s){var o=e+(t&r|~t&a)+(n>>>0)+s;return(o<<i|o>>>32-i)+t},s._gg=function(e,t,r,a,n,i,s){var o=e+(t&a|r&~a)+(n>>>0)+s;return(o<<i|o>>>32-i)+t},s._hh=function(e,t,r,a,n,i,s){var o=e+(t^r^a)+(n>>>0)+s;return(o<<i|o>>>32-i)+t},s._ii=function(e,t,r,a,n,i,s){var o=e+(r^(t|~a))+(n>>>0)+s;return(o<<i|o>>>32-i)+t},s._blocksize=16,s._digestsize=16,e.exports=function(e,r){if(void 0===e||null===e)throw new Error("Illegal argument "+e);var a=t.wordsToBytes(s(e,r));return r&&r.asBytes?a:r&&r.asString?i.bytesToString(a):t.bytesToHex(a)}})()},"9a634":function(e,t){var r={utf8:{stringToBytes:function(e){return r.bin.stringToBytes(unescape(encodeURIComponent(e)))},bytesToString:function(e){return decodeURIComponent(escape(r.bin.bytesToString(e)))}},bin:{stringToBytes:function(e){for(var t=[],r=0;r<e.length;r++)t.push(255&e.charCodeAt(r));return t},bytesToString:function(e){for(var t=[],r=0;r<e.length;r++)t.push(String.fromCharCode(e[r]));return t.join("")}}};e.exports=r},ac2a:function(e,t,r){"use strict";r.r(t);var a=function(){var e=this,t=e.$createElement,r=e._self._c||t;return r("div",{staticClass:"main"},[r("a-form",{ref:"formLogin",staticClass:"user-layout-login",attrs:{id:"formLogin",form:e.form},on:{submit:e.handleSubmit}},[r("a-tabs",{attrs:{activeKey:e.customActiveKey,tabBarStyle:{textAlign:"center",borderBottom:"unset"}},on:{change:e.handleTabClick}},[r("a-tab-pane",{key:"tab1",attrs:{tab:e.$t("user.login.tab-login-credentials")}},[e.isLoginError?r("a-alert",{staticStyle:{"margin-bottom":"24px"},attrs:{type:"error",showIcon:"",message:e.userloginmessage}}):e._e(),r("a-form-item",[r("a-input",{directives:[{name:"decorator",rawName:"v-decorator",value:["username",{rules:[{required:!0,message:e.$t("user.userName.required")},{validator:e.handleUsernameOrEmail}],validateTrigger:"change"}],expression:"[\n              'username',\n              {rules: [{ required: true, message: $t('user.userName.required') }, { validator: handleUsernameOrEmail }], validateTrigger: 'change'}\n            ]"}],attrs:{size:"large",type:"text",placeholder:e.$t("user.login.username.placeholder")}},[r("a-icon",{style:{color:"rgba(0,0,0,.25)"},attrs:{slot:"prefix",type:"user"},slot:"prefix"})],1)],1),r("a-form-item",[r("a-input-password",{directives:[{name:"decorator",rawName:"v-decorator",value:["password",{rules:[{required:!0,message:e.$t("user.password.required")}],validateTrigger:"blur"}],expression:"[\n              'password',\n              {rules: [{ required: true, message: $t('user.password.required') }], validateTrigger: 'blur'}\n            ]"}],attrs:{size:"large",placeholder:e.$t("user.login.password.placeholder")}},[r("a-icon",{style:{color:"rgba(0,0,0,.25)"},attrs:{slot:"prefix",type:"lock"},slot:"prefix"})],1)],1)],1),r("a-tab-pane",{key:"tab2",attrs:{tab:e.$t("user.login.tab-login-mobile")}},[r("a-form-item",[r("a-input",{directives:[{name:"decorator",rawName:"v-decorator",value:["mobile",{rules:[{required:!0,pattern:/^1[34578]\d{9}$/,message:e.$t("user.login.mobile.placeholder")}],validateTrigger:"change"}],expression:"['mobile', {rules: [{ required: true, pattern: /^1[34578]\\d{9}$/, message: $t('user.login.mobile.placeholder') }], validateTrigger: 'change'}]"}],attrs:{size:"large",type:"text",placeholder:e.$t("user.login.mobile.placeholder")}},[r("a-icon",{style:{color:"rgba(0,0,0,.25)"},attrs:{slot:"prefix",type:"mobile"},slot:"prefix"})],1)],1),r("a-row",{attrs:{gutter:16}},[r("a-col",{staticClass:"gutter-row",attrs:{span:16}},[r("a-form-item",[r("a-input",{directives:[{name:"decorator",rawName:"v-decorator",value:["captcha",{rules:[{required:!0,message:e.$t("user.verification-code.required")}],validateTrigger:"blur"}],expression:"['captcha', {rules: [{ required: true, message: $t('user.verification-code.required') }], validateTrigger: 'blur'}]"}],attrs:{size:"large",type:"text",placeholder:e.$t("user.login.mobile.verification-code.placeholder")}},[r("a-icon",{style:{color:"rgba(0,0,0,.25)"},attrs:{slot:"prefix",type:"mail"},slot:"prefix"})],1)],1)],1),r("a-col",{staticClass:"gutter-row",attrs:{span:8}},[r("a-button",{staticClass:"getCaptcha",attrs:{tabindex:"-1",disabled:e.state.smsSendBtn},domProps:{textContent:e._s(!e.state.smsSendBtn&&e.$t("user.register.get-verification-code")||e.state.time+" s")},on:{click:function(t){return t.stopPropagation(),t.preventDefault(),e.getCaptcha(t)}}})],1)],1)],1)],1),r("a-form-item",[r("a-checkbox",{directives:[{name:"decorator",rawName:"v-decorator",value:["rememberMe",{valuePropName:"checked"}],expression:"['rememberMe', { valuePropName: 'checked' }]"}]},[e._v(e._s(e.$t("user.login.remember-me")))]),r("router-link",{staticClass:"forge-password",staticStyle:{float:"right"},attrs:{to:{name:"recover",params:{user:"aaa"}}}},[e._v(e._s(e.$t("user.login.forgot-password")))])],1),r("a-form-item",{staticStyle:{"margin-top":"24px"}},[r("a-button",{staticClass:"login-button",attrs:{size:"large",type:"primary",htmlType:"submit",loading:e.state.loginBtn,disabled:e.state.loginBtn}},[e._v(e._s(e.$t("user.login.login")))])],1),r("div",{staticClass:"user-login-other"},[r("span",[e._v(e._s(e.$t("user.login.sign-in-with")))]),r("a",[r("a-icon",{staticClass:"item-icon",attrs:{type:"alipay-circle"}})],1),r("a",[r("a-icon",{staticClass:"item-icon",attrs:{type:"taobao-circle"}})],1),r("a",[r("a-icon",{staticClass:"item-icon",attrs:{type:"weibo-circle"}})],1),r("router-link",{staticClass:"register",attrs:{to:{name:"register"}}},[e._v(e._s(e.$t("user.login.signup")))])],1)],1)],1)},n=[],i=r("5530"),s=(r("d3b7"),r("6821")),o=r.n(s),c=r("5880"),u=r("ca00"),l=r("7ded"),g={components:{},data:function(){return{customActiveKey:"tab1",loginBtn:!1,loginType:0,isLoginError:!1,requiredTwoStepCaptcha:!0,stepCaptchaVisible:!0,form:this.$form.createForm(this),state:{time:60,loginBtn:!1,loginType:0,smsSendBtn:!1},userloginmessage:""}},created:function(){},methods:Object(i["a"])(Object(i["a"])({},Object(c["mapActions"])(["Login","Logout"])),{},{handleUsernameOrEmail:function(e,t,r){var a=this.state,n=/^([a-zA-Z0-9_-])+@([a-zA-Z0-9_-])+((\.[a-zA-Z0-9_-]{2,3}){1,2})$/;n.test(t)?a.loginType=0:a.loginType=1,r()},handleTabClick:function(e){this.customActiveKey=e},handleSubmit:function(e){var t=this;e.preventDefault();var r=this.form.validateFields,a=this.state,n=this.customActiveKey,s=this.Login;a.loginBtn=!0;var c="tab1"===n?["username","password"]:["mobile","captcha"];r(c,{force:!0},(function(e,r){if(e)setTimeout((function(){a.loginBtn=!1}),600);else{var n=Object(i["a"])({},r);delete n.username,n[a.loginType?"username":"email"]=r.username,n.password=o()(r.password),s(n).then((function(e){return t.loginSuccess(e)})).catch((function(e){return t.requestFailed(e)})).finally((function(){a.loginBtn=!1}))}}))},getCaptcha:function(e){var t=this;e.preventDefault();var r=this.form.validateFields,a=this.state;r(["mobile"],{force:!0},(function(e,r){if(!e){a.smsSendBtn=!0;var n=window.setInterval((function(){a.time--<=0&&(a.time=60,a.smsSendBtn=!1,window.clearInterval(n))}),1e3),i=t.$message.loading("验证码发送中..",0);Object(l["b"])({mobile:r.mobile}).then((function(e){setTimeout(i,2500),t.$notification["success"]({message:"提示",description:"验证码获取成功，您的验证码为："+e.result.captcha,duration:8})})).catch((function(e){setTimeout(i,1),clearInterval(n),a.time=60,a.smsSendBtn=!1,t.requestFailed(e)}))}}))},stepCaptchaSuccess:function(){this.loginSuccess()},stepCaptchaCancel:function(){var e=this;this.Logout().then((function(){e.loginBtn=!1,e.stepCaptchaVisible=!1}))},loginSuccess:function(e){var t=this;this.$router.push({path:"/"}),setTimeout((function(){t.$notification.success({message:"欢迎",description:"".concat(Object(u["b"])(),"，欢迎回来")})}),1e3),this.isLoginError=!1},requestFailed:function(e){this.isLoginError=!0,this.userloginmessage=e.msg||"请求出现错误，请稍后再试"}})},d=g,f=(r("e649"),r("2877")),m=Object(f["a"])(d,a,n,!1,null,"62715b39",null);t["default"]=m.exports},e649:function(e,t,r){"use strict";r("1608")}}]);