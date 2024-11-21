(function (factory) {
  if (typeof define === 'function' && define.amd) {
      // AMD. Register as an anonymous module.
      define(['jquery'], factory);
  } else if (typeof exports === 'object') {
      // Node/CommonJS style for Browserify
      module.exports = factory;
  } else {
      // Browser globals
      factory(jQuery);
  }
}(function($){
  var cachedLang=null;
	var htmlEncodeRegexp=/&|<|>| |\'|\"/g,htmlEncodeMapping = {
		"&":'&amp;',
		"<":'&lt;',
		">":'&gt;',
		" ":'&nbsp;',
		"'":'&#39;',
		'"':'&quot;',
	}
	var htmlDecodeRegexp=/&amp;|&lt;|&gt;|&nbsp;|&#39;|&quot;/g,htmlDecodeMapping = {
		'&amp;':"&",
		'&lt;':"<",
		'&gt;':">",
		'&nbsp;':" ",
		'&#39;':"'",
		'&quot;':'"',
	}
  var textNl2brRegexp=/\n|  |\t/g,textNl2brMapping={
    "\n":'<br />',
    "  ":'&nbsp; ',
    "\t":'&nbsp; &nbsp; ',
  }
	function getJQueryObject(a){
		return (typeof(a)=='object' && a instanceof jQuery) ? a : $(a);
	}
  window.App={
    clientID: {},
    i18n: {
			SYS_INFO: 'System Information', 
			UPLOAD_ERR: 'Upload Error', 
			PLEASE_SELECT_FOR_OPERATE: 'Please select the item you want to operate', 
			PLEASE_SELECT_FOR_REMOVE: 'Please select the item you want to delete', 
			CONFIRM_REMOVE: 'Are you sure you want to delete them?', 
			SELECTED_ITEMS: 'You have selected %d items', 
			SUCCESS: 'The operation was successful', 
			FAILURE: 'Operation failed', 
			UPLOADING:'File uploading, please wait...', 
			UPLOAD_SUCCEED:'Upload successfully', 
			BUTTON_UPLOAD:'Upload' 
		},
    lang: 'en',
    status: {CaptchaError:-9,CaptchaIdMissing:-10,CaptchaCodeRequired:-11,NonPrivileged:-2,NotLoggedIn:-1,Failure:0,Success:1},
    sprintf:sprintfWrapper.init,
    t: function(key) {
      if(typeof(App.i18n[key])=='undefined'){
        if(arguments.length < 2) return key;
        return App.sprintf.apply(this,arguments);
      }
      if(arguments.length < 2) return App.i18n[key];
      arguments[0]=App.i18n[key];
      return App.sprintf.apply(this,arguments);
    },
    langInfo: function(){
      if(cachedLang!=null) return cachedLang;
      var _lang=App.lang.split('-',2);
      cachedLang={encoding:_lang[0],country:''};
      if(_lang.length>1)cachedLang.country=_lang[1].toUpperCase();
      return cachedLang;
    },
    langTag: function (seperator) {
      var l = App.langInfo();
      if (l.country) {
        if(seperator==null) seperator = '-';
        return l.encoding + seperator + l.country;
      }
      return l.encoding;
    },
		htmlEncode: function(value){
			if(!value) return value;
			return String(value).replace(htmlEncodeRegexp, function(v){
				return htmlEncodeMapping[v];
			});
		},
		getJQueryObject: getJQueryObject,
		htmlDecode: function(value){
			if(!value) return value;
			return String(value).replace(htmlDecodeRegexp, function(v){
				return htmlDecodeMapping[v];
			});
		},
    dialog:function(){
      //requirement: https://github.com/aafrontend/bootstrap3-dialog
      return BootstrapDialog;
    },
    message:function(options,sticky){
			if (typeof (options) == 'string') {
				switch(options){
					case 'remove':
						var number = sticky;
						return $.gritter.remove(number);
					case 'clear': //$.gritter.removeAll({before_close:function(wrap){},after_close:function(){}});
						return $.gritter.removeAll(sticky||{});
				}
			}
      var defaults={title: App.i18n.SYS_INFO,
        text: '',
        image: '',
        class_name: 'clean',//primary|info|danger|warning|success|dark
        sticky: false
        //,time: 1000,speed: 500,position: 'bottom-right'
      };
      if(typeof(options)!="object")options={text:options,type:'clean'};
      if(('type' in options)&&options.type)options.class_name=options.type;
      options=$.extend({},defaults,options||{});
      switch(options.class_name){
        case 'dark':
        case 'primary':
        case 'clean':
        case 'info':
        if(options.title) options.title='<i class="fa fa-info-circle"></i> '+App.t(options.title);break;

        case 'error':
        options.class_name='danger';
        case 'danger':
        if(options.title) options.title='<i class="fa fa-comment-o"></i> '+App.t(options.title);break;

        case 'warning':
        if(options.title) options.title='<i class="fa fa-warning"></i> '+App.t(options.title);break;
          
        case 'success':
        if(options.title) options.title='<i class="fa fa-check"></i> '+App.t(options.title);break;
      }
      if(sticky!=null)options.sticky=sticky;
      if(options.text)options.text=App.t(options.text);
      var number = $.gritter.add(options);
      return number;
    },
		text2html: function (text, noescape) {
			text = String(text);
			if(!noescape) text.replace(/<|>/g, function(v){
				return v=='<'?'&lt;':'&gt;';
			});
			return App.textNl2br(text);
		},
		ifTextNl2br: function (text) {
			text = String(text);
			if (/<[^>]+>/.test(text)) return text;
			return App.textNl2br(text);
		},
		textNl2br: function (text) {
			return text.replace(textNl2brRegexp, function(v){
				return textNl2brMapping[v];
			});
		},
		trimSpace: function (text) {
			return String(text).replace(/^[\s]+|[\s]+$/g,'');
		},
    checkedAll: function (ctrl, target) {
      return $(target).not(':disabled').prop('checked', $(ctrl).prop('checked'));
    },
    attachCheckedAll: function (ctrl, target, showNumElem) {
      $(ctrl).on('click', function () {
        App.checkedAll(this, target);
        if (showNumElem) $(showNumElem).text($(target + ':checked').length);
      });
    },
    attachAjaxURL:function(elem){
			if (elem == null) elem = document;
			$(elem).on('click', '[data-ajax-url]', function () {
				var a = $(this), confirmMsg = a.data('ajax-confirm');
				if(a.data('processing')){
					alert(App.t('Processing, please wait for the operation to complete'));
					return;
				}
				if(confirmMsg && !confirm(confirmMsg)) return;
				a.data('processing',true);
				var url = a.data('ajax-url'), method = a.data('ajax-method') || 'get', params = a.data('ajax-params') || {}, title = a.attr('title')||App.i18n.SYS_INFO, accept = a.data('ajax-accept') || 'html', target = a.data('ajax-target'), callback = a.data('ajax-callback'), toggle = a.data('ajax-toggle'), onsuccess = a.data('ajax-onsuccess'), reload = a.data('ajax-reload') || false;
				if (!title) title = a.text();
				var fa = a.children('.fa');
				var hasIcon = toggle && fa.length>0;
				if (hasIcon){
					fa.addClass('fa-spin')
				}else{
					App.loading('show');
				}
				a.trigger('processing');
				if (typeof params === "function") params = params.apply(this, arguments);
				//params = App.setClientID(params);
				$[method](url, params || {}, function (r) {
					a.data('processing',false);
					a.trigger('finished',arguments);
					if (hasIcon){
						fa.removeClass('fa-spin');
					}else{
						App.loading('hide');
					}
					if (callback) return callback.apply(this, arguments);
					if (target) {
						var data;
						if (accept == 'json') {
							if (r.Code != 1) {
								return App.message({ title: title, text: r.Info, type: 'error', time: 5000, sticky: false });
							}
							data = r.Data;
						} else {
							data = r;
						}
						$(target).html(data);
						a.trigger('partial.loaded', arguments);
						$(target).trigger('partial.loaded', arguments);
						if(onsuccess) window.setTimeout(onsuccess,0);
						return;
					}
					if(r.Code == 1){
						if(onsuccess) window.setTimeout(onsuccess,2000);
						if(reload) window.setTimeout(function(){window.location.reload()},2000);
					}
					if (accept == 'json') {
						return App.message({ title: title, text: r.Info, type: r.Code == 1 ? 'success' : 'error', time: 5000, sticky: false });
					}
					App.message({ title: title, text: r, time: 5000, sticky: false });
				}, accept).error(function (xhr, status, info) {
					a.data('processing',false);
					a.trigger('finished',arguments);
					if (hasIcon){
						fa.removeClass('fa-spin');
					}else{
						App.loading('hide');
					}
					App.message({ title: title, text: xhr.responseText, type: 'error', time: 5000, sticky: false });
				});
			});
    },
    attachPjax:function(elem,callbacks,timeout){
      if(!$.support.pjax)return;
      if(elem==null)elem='a';
      if(timeout==null)timeout=5000;
      var defaults={onclick:null,onsend:null,oncomplete:null,ontimeout:null,onstart:null,onend:null};
      var options=$.extend({},defaults,callbacks||{});
      $(document).on('click', elem+'[data-pjax]', function(event) {
        var container = $(this).data('pjax'),keepjs=$(this).data('keepjs');
        var onclick=$(this).data('onclick');
        $.pjax.click(event, $(container),{timeout:timeout,keepjs:keepjs});
        if(options.onclick)options.onclick(this);
        if(onclick && typeof(window[onclick])=='function')window[onclick](this);
      }).on('pjax:send',function(evt,xhr,option){
        App.loading('show');
        if(options.onsend)options.onsend(evt,xhr,option);
      }).on('pjax:complete',function(evt, xhr, textStatus, option){
        App.loading('hide');
        if(options.oncomplete)options.oncomplete(evt, xhr, textStatus, option);
      }).on('pjax:timeout',function(evt,xhr,option){
        console.log('timeout');
        App.loading('hide');
        if(options.ontimeout)options.ontimeout(evt,xhr,option);
      }).on('pjax:start',function(evt,xhr,option){
        if(options.onstart)options.onstart(evt,xhr,option);
      }).on('pjax:end',function(evt,xhr,option){
        App.loading('hide');
        if(options.onend)options.onend(evt,xhr,option);
        //console.debug(option);
        var id=option.container.attr('id');
        if(id){
          App.bottomFloat('#'+id+' .pagination');
          App.bottomFloat('#'+id+' .form-submit-group',0,true);
          $('#'+id+' .switch:not(.has-switch)').bootstrapSwitch();
        }
        if(option.type=='GET') $('#global-search-form').attr('action',option.url);
      });
    },
    loading:function(op,n){
      switch(op){
        case 'show':
            NProgress.start();break;
        case 'hide':
            NProgress.done();break;
        case 'set':
            NProgress.set(n);break;
        case 'inc':
            NProgress.inc(n);break;
      }
    },
    insertAtCursor: function(myField, myValue,posStart,posEnd) { 
      if (typeof TextAreaEditor != 'undefined') {
        TextAreaEditor.setSelectText(myField, myValue,posStart,posEnd);
        return;
      } 
      /* IE support */
			if (document.selection) {
				myField.focus();
				sel = document.selection.createRange();
				sel.text = myValue;
				sel.select();
			} /* MOZILLA/NETSCAPE support */
			else if (myField.selectionStart || myField.selectionStart == '0') {
				var startPos = myField.selectionStart;
				var endPos = myField.selectionEnd; /* save scrollTop before insert */
				var restoreTop = myField.scrollTop;
				myField.value = myField.value.substring(0, startPos) + myValue + myField.value.substring(endPos, myField.value.length);
				if (restoreTop > 0) myField.scrollTop = restoreTop;
				myField.focus();
				myField.selectionStart = startPos + myValue.length;
				myField.selectionEnd = startPos + myValue.length;
			} else {
				myField.value += myValue;
				myField.focus();
      }
    },
    wsURL:function(url) {
    	var protocol='ws:';
      if(window.location.protocol=='https:')protocol='wss:';
      var p=String(url).indexOf('//');
      if(p==-1){
        url=protocol+"//"+window.location.host+url;
      }else{
        url=protocol+String(url).substring(p);
      }
      return url;
    },
    websocket:function(showmsg,url,onopen,onclose){
      if(!('WebSocket' in window)) return;
      url = App.wsURL(url);
    	var ws = new WebSocket(url);
    	ws.onopen = function(evt) {
    	    console.log('Websocket Server is connected');
    		  if(onopen!=null&&$.isFunction(onopen))onopen.apply(this,arguments);
    	};
    	ws.onclose = function(evt) {
    	    console.log('Websocket Server is disconnected');
          if (onclose != null && typeof onclose === "function") onclose.apply(this, arguments);
    	};
    	ws.onmessage = function(evt) {
    	    if(showmsg) showmsg(evt.data);
    	};
    	//ws.onerror = function(evt) {console.dir(evt);};
      if(onopen!=null&&typeof(onopen)=='object'){
        ws=$.extend({},ws,onopen);
      }
      return ws;
    },
    notifyListen:function(notifyURL){
      if(notifyURL==null) notifyURL='/user/notice';
      var messageCount={notify:0,element:0,modal:0},
      messageMax={notify:20,element:50,modal:50};
      App.websocket(function(message){
        //console.dir(message);
        var m=$.parseJSON(message);
        if(typeof(App.clientID['notify'])=='undefined'){
          App.clientID['notify']=m.client_id;
        }
        switch(m.mode){
          case '-':
          break;

          case 'element':
          var c=$('#notify-element-'+m.type);
          if(c.length<1){
            var callback='recv_notice_'+m.type;
            if(typeof(window[callback])!='undefined'){
              return window[callback](m);
            }
            if(m.status>0){
              console.info(m.content);
            }else{
              console.error(m.content);
            }
            return;
          }
          if(messageCount[m.mode]>=messageMax[m.mode]){
            c.find('li:first').remove();
          }
          if(m.title){
            var badge='badge-danger';
            if(m.status>0) badge='badge-success';
            message='<span class="badge '+badge+'">'+App.text2html(m.title)+'</span> '+App.text2html(m.content);
          }else{
            message=App.text2html(m.content);
          }
          c.append('<li>'+message+'</li>');
          messageCount[m.mode]++;
          break;

          case 'modal':
          var c=$('#notify-modal-'+m.type);
          if(c.length<1){
            var callback='recv_notice_'+m.type;
            if(typeof(window[callback])!='undefined'){
              return window[callback](m);
            }
            if(m.status>0){
              console.info(m.content);
            }else{
              console.error(m.content);
            }
            return;
          }
          if(m.title){
            var badge='badge-danger';
            if(m.status>0) badge='badge-success';
            message='<span class="badge '+badge+'">'+App.text2html(m.title)+'</span> '+App.text2html(m.content);
          }else{
            message=App.text2html(m.content);
          }
          if(!c.data('shown')){
            messageCount[m.mode]=0;
            var mbody=c.find('.modal-body'),mbodyUL=mbody.children('ul.modal-body-ul');
            if(mbodyUL.length<1){
              mbody.html('<ul class="modal-body-ul" id="notify-modal-'+m.type+'-container"><li>'+message+'</li></ul>');
            }else{
              mbodyUL.html('<li>'+message+'</li>');
            }
            App.dialog().show({
              id:'notify-modal-'+m.type,
              title:App.t('身份验证'),
              type:'type-primary',
              message:html,
              nl2br:false,
              closeByBackdrop:false,
              onshown: function(dialogRef){
                c.data('shown',true);
              },
              onhidden: function(dialogRef){
                c.data('shown',false);
              }
            });
          }else{
            var cc=$('#notify-modal-'+m.type+'-container');
            if(messageCount[m.mode]>=messageMax[m.mode]){
              cc.find('li:first').remove();
            }
            cc.append('<li>'+message+'</li>');
          }
          messageCount[m.mode]++;
          break;
          
          case 'notify':
          default:
            if('notify'!=m.mode) m.mode='notify';
            var c=$('#notice-message-container');
            if(c.length<1){
              App.message({title: App.i18n.SYS_INFO, text: '<ul id="notice-message-container" class="no-list-style" style="max-height:500px;overflow-y:auto;overflow-x:hidden"></ul>',sticky: true});
              c=$('#notice-message-container');
            }
            if(messageCount[m.mode]>=messageMax[m.mode]){
              c.find('li:first').remove();
            }
            if(m.title){
              var badge='badge-danger';
              if(m.status>0) badge='badge-success';
              message='<span class="badge '+badge+'">'+App.text2html(m.title)+'</span>'+App.text2html(m.content);
            }else{
              message=App.text2html(m.content);
            }
            c.append('<li>'+message+'</li>');
            messageCount[m.mode]++;
          break;
        }
      },notifyURL);
    },
    bottomFloat:function (elems, top, autoWith) {
        if ($(elems).length<1) return;
        if (top == null) top = 0;
        $(elems).not('[disabled-fixed]').each(function(){
        var elem=this;
      	var _offset = $(elem).height() + top;
      	var offsetY = $(elem).offset().top + _offset;
      	var w = $(elem).outerWidth(), h = $(elem).outerHeight();
      	if (autoWith||$(elem).data('auto-width')) $(elem).css('width', w);
      	$(window).on('scroll',function () {
      		var scrollH = $(this).scrollTop() + $(window).height();
      		if (scrollH >= offsetY) {
      			if ($(elem).hasClass('always-bottom')) {
      				$(elem).removeClass('always-bottom');
      				$(elem).next('.fixed-placeholder').hide();
      			}
      			return;
      		}
      		if (!$(elem).hasClass('always-bottom')) {
      			$(elem).addClass('always-bottom');
      			if ($(elem).next('.fixed-placeholder').length > 0) {
      				$(elem).next('.fixed-placeholder').show();
      			} else {
      				$(elem).after('<div style="width:' + w + 'px;height:' + h + 'px" class="fixed-placeholder"></div>');
      			}
      		}
      	});
      });
      $(window).trigger('scroll');
    },
    getImgNaturalDimensions:function (oImg, callback) {
      if (oImg.naturalWidth) { // 现代浏览器
        callback({w: oImg.naturalWidth, h:oImg.naturalHeight});
        return;
      } 
      // IE6/7/8
      var nImg = new Image();
      nImg.onload = function() {
        callback({w: nImg.width, h:nImg.height});
      }
      nImg.src = oImg.src;
    },
    replaceURLParam:function(name,value,url){
      if(url==null) url=window.location.href;
      value=encodeURIComponent(value);
      var pos=String(url).indexOf('?');
      if(pos<0) return url+'?'+name+'='+value;
      var q=url.substring(pos),r=new RegExp('([\\?&]'+name+'=)[^&]*(&|$)');
      if(!r.test(q)) return url+'&'+name+'='+value;
      url=url.substring(0,pos);
	    q=q.replace(r,'$1'+value+'$2');
	    return url+q;
    },
    switchLang:function(lang){
      window.location=App.replaceURLParam('lang',lang);
    },
    tableSorting: function (table) {
      table = table == null ? '' : table + ' ';
      $(table + '[sort-current!=""]').each(function () {//<thead sort-current="created">
        var current = String($(this).attr('sort-current'));
        var isDesc = current.substring(0,1) == '-';
        if (isDesc) current=current.substring(1);
        var sortObj = $(this).find('[sort="' + current + '"]');//<th sort="-created">
        var newCls, oldCls, sortBy;
        if (sortObj.length < 1 && current) {
          sortObj = $(this).find('[sort="-' + current + '"]');
        }
        if (!isDesc) {
          newCls = 'fa-arrow-up';
          sortBy = 'up';
          oldCls = 'fa-arrow-down';
        } else {
          newCls = 'fa-arrow-down';
          sortBy = 'down';
          oldCls = 'fa-arrow-up';
        }
        if (sortObj.length > 0) {
          var icon = sortObj.children('.fa');
          if (icon.length < 1) {
            sortObj.append('<i class="fa ' + newCls + '"></i>');
          } else {
            icon.removeClass(oldCls).addClass(newCls);
          }
          sortObj.addClass('sort-active sort-'+sortBy);
          sortObj.siblings('.sort-active').removeClass('sort-active').removeClass('sort-up').removeClass('sort-down');
        }
      });
      $(table + '[sort-current] [sort]').css('cursor', 'pointer').on('click', function (e) {
        var thead = $(this).parents('[sort-current]');
        var current = thead.attr('sort-current');
        var url = thead.attr('sort-url') || window.location.href;
        var trigger = thead.attr('sort-trigger');
        var sort = $(this).attr('sort');
        if (current && (current == sort || current == '-' + sort)) {
          var reg = /^\-/;
          current = reg.test(current) ? current.replace(reg, '') : '-' + current;
        } else {
          current = sort;
        }
        url = App.replaceURLParam('sort', current, url);
        if (trigger) {
          thead.trigger('sort');
        } else {
          var setto = thead.attr('sort-setto');
          if (setto) {
            $(setto).load(url);
          } else {
            window.location = url;
          }
        }
      });
    },
    formatBytes:function(bytes,precision) {
      if(precision==null)precision=2;
      var units = ["YB", "ZB", "EB", "PB", "TB", "GB", "MB", "KB", "B"];
      var total=units.length;
      for(total--; total > 0 && bytes > 1024.0; total--) {
        bytes /= 1024.0;
      }
      return bytes.toFixed(precision)+units[total];
    },
    progressMonitor:function (getCurrentFn,totalProgress){
      NProgress.start();
      var interval=window.setInterval(function(){
        var current=getCurrentFn()/totalProgress;
        if(current>=1){
          NProgress.set(1);
          window.clearInterval(interval);
        }else{
          NProgress.set(current);
        }
      },50);
    },
    float:function (elem, mode, attr, position, options) {
      if (!mode) mode = 'ajax'; 
      if (!attr) attr = mode=='remind'?'rel':'src'; 
			if (!position) position = '5-7';//两个数字分别代表trigger(浮动层)-target(来源对象)，（各个数字的编号从矩形框的左上角开始，沿着顺时针开始旋转来进行编号，然后再从上中部开始沿着顺时针开始编号进行。也就是1、2、3、4分别代表左上角、右上角、右下角、左下角；5、6、7、8分别代表上中、右中、下中、左中）
			else {
				switch (position) {
					case 'bottom':position='5-7';break;
					case 'right':position='8-6';break;
					case 'top':position='7-5';break;
					case 'left':position='6-8';break;
					case 'left-bottom':position='2-4';break;
					case 'right-bottom':position='1-3';break;
					case 'left-top':position='3-1';break;
					case 'right-top':position='4-2';break;
				}
			}
      var defaults = { 'targetMode': mode, 'targetAttr': attr, 'position': position };
      $(elem).powerFloat($.extend(defaults,options||{}));
    },
		uploadPreviewer: function (elem, options, callback) {
			if($(elem).parent('.file-preview-shadow').length<1){
				var defaults = {
					"buttonText":'<i class="fa fa-cloud-upload"></i> '+App.i18n.BUTTON_UPLOAD,
					"previewTableContainer":'#previewTableContainer',
					"url":'',
					"previewTableShow":false,
					"uploadProgress":function(progress){
						var count=progress*100;
						if(count>100){
							$.LoadingOverlay("hide");
							return;
						}
						$.LoadingOverlay("progress", count);
					}
				};
				var noptions = $.extend({}, defaults, options || {});
				var uploadInput = $(elem).uploadPreviewer(noptions);
				$(elem).data('uploadPreviewer', uploadInput);
				$(elem).on("file-preview:changed", function(e) {
					var options = {
						image : ASSETS_URL+"/images/nging-gear.png", 
						//fontawesome : "fa fa-cog fa-spin",
						text  : App.i18n.UPLOADING
					};
					if(noptions.uploadProgress){
						options.progress = true;
						options.image = "";
					}
				  	$.LoadingOverlay("show", options);
				  	uploadInput.submit(function(r){
					  $.LoadingOverlay("hide");
					  if(r.Code==1){
						  App.message({text:App.i18n.UPLOAD_SUCCEED,type:'success'});
					  }else{
						  App.message({text:r.Info,type:'error'});
					  }
					  if(callback!=null) callback.call(this, r);
				  	});
				});
			}
		},
    showRequriedInputStar:function () {
      $('form:not([required-redstar])').each(function(){
        $(this).find('[required]').each(function(){
          var id = $(this).attr('id');
          if($('label[for=' + id + ']').length > 0){
            $('label[for=' + id + ']').addClass('required');
						return;
          }
					var parent = $(this).parent('.input-group');
					if(parent.length>0){
						parent.addClass('required');
						return;
					}
					parent = $(this).parent('div[class*="col-"]');
					if (parent.length>0) {
            if(parent.prev('.control-label').length>0){
              parent.prev('.control-label').addClass('required');
              return;
            }
            if(parent.prev('.col-form-label').length>0){
              parent.prev('.col-form-label').addClass('required');
              return;
            }
					}
          var row = $(this).closest('.form-group');
          if (row.length<1) return;
          var lbl = row.children('.control-label:not(.required),.col-form-label:not(.required)');
          if (lbl.length<1) return;
          lbl.addClass('required');
        });
				$(this).attr('required-redstar','1');
      });
    },
		pushState:function(data,title,url){
			if(!window.history || !window.history.pushState)return;
			window.history.pushState(data,title,url);
		},
		replaceState:function(data,title,url){
			if(!window.history || !window.history.replaceState)return;
			window.history.replaceState(data,title,url);
		},
		formatJSON:function(json){
			json = $.trim(json);
			var first = json.substring(0,1);
			if (first=='['||first=='{'){
				var obj = JSON.parse(json);
				json = JSON.stringify(obj, null, "\t");
				return json;
			}
			return '';
		},
		formatJSONFromInnerHTML:function($a){
			if($a.data('jsonformatted'))return;
			$a.data('jsonformatted',true);
			var json = App.formatJSON($a.html());
			if(json!='') $a.html(json);
		},
		captchaUpdate: function($form, resp){
			if(!App.captchaHasError(resp.Code) || !resp.Data || typeof(resp.Data.captchaType) === 'undefined') return;
			if(false == ($form instanceof jQuery)) $form=$($form);
      switch(resp.Data.captchaType){
        case 'api':
          if(resp.Data.jsURL && $('script[src="'+resp.Data.jsURL+'"]').length<1){
            $('body').append('<script src="'+resp.Data.jsURL+'" type="text/javascript"></script>');
            if(resp.Data.jsInit){
              eval(resp.Data.jsInit);
            }
            return;
          }
          if(resp.Data.jsCallback) {
            eval('('+resp.Data.jsCallback+')();');
          }
          break;
        default:
          var idElem = $form.find('input#'+resp.Data.captchaIdent);
          idElem.val(resp.Data.captchaID);
          idElem.siblings('img').attr('src',resp.Data.captchaURL);
          if(resp.Data.captchaName) $form.find('input[name="'+resp.Data.captchaName+'"]').focus();
      }
		},
		captchaHasError: function(code) {
			return code >= -11 && code <= -9;
		},
    fixedFooter: function(elem,offset) {
      var $footer=$(elem);
      if($footer.length<1)return;
      var marginTop=$footer.css("marginTop");
      if(offset==null) offset = $('.pos-fixed-top:first').height();
      var fixer=function($footer){
        var b=$footer.position().top+$footer.height()+offset;
        var f=$(window).height()-b;
        if(f>0) $footer.css('margin-top',f+'px');
        else if(f<0) $footer.css('margin-top',marginTop);
      }
      fixer($footer);
      $(window).on('resize',function(){
        fixer($footer);
      });
    }
};
return window.App;
}));