/**
 * 显示对话框
 * @param {string} msg 
 * @param {string} type 
 */
function showModal(msg,type){
    var oldClass=$('#error-dialog').data('old-class')||'danger';
    var btn=$('#error-dialog').find('#error-button[data-dismiss="modal"]');
    var cls=btn.attr('class');
    var match=String(cls).match(/btn-([^ ]+)/);
    var icon,oldIcon=$('#error-dialog').data('old-icon')||'ion-ios-close-outline';
    $('#error-dialog').data('old-class',match[1]);
    if(type==null) type='error';
    switch(type){
        case 'error':
        type='danger';
        icon='ion-ios-close-outline';
        break;
        case 'success':
        icon='ion-ios-checkmark-outline';
        break;
        default:
        type='info';
        icon='ion-info';
        break;
    }
    $('#error-dialog').data('old-icon',icon);
    if(oldClass){
        btn.removeClass('btn-'+oldClass).addClass('btn-'+type);
        $('#error-dialog').find('.tx-'+oldClass).removeClass('tx-'+oldClass).addClass('tx-'+type);
        $('#error-dialog').find('.'+oldIcon).removeClass(oldIcon).addClass(icon);
    }
    $('#error-dialog').find('#error-message').text(msg);
    $('#error-dialog').modal('show');
}
/**
 * 隐藏对话框
 */
function hideModal(){
    $('#error-dialog').modal('hide');
}
/**
 * ajax登录
 * @param {string|object} elem 
 * @param {string} nextURL 
 */
function signIn(elem,nextURL){
    var cc=checkSubmitBtnWithCloser(elem);
    if(cc.submited)return;
    if(nextURL==null) nextURL=window.location.href;
    var data=$(elem).serializeArray();
    var url=$(elem).attr('action')||BASE_URL+'/sign_in';
    var ajaxOptions={
        url: url,
        type: 'POST',
        data: data,
        dataType: 'json',
        close: cc.close,
        success: function(r){
            if(r.Code!=1) {
                if(App.captchaHasError(r.Code) && r.Data && typeof(r.Data.captchaName)!='undefined' && r.Data.captchaName && $(elem).find('input[name="'+r.Data.captchaName+'"]').length<1){ 
                    return captchaDialog(r,ajaxOptions);
                }
                cc.close()
                showMsg({text:r.Info,type:'error'});
                return renewCaptcha(elem,r);
            }
            cc.close()
            showMsg({text:r.Info,type:'success'});
            var callback=$(elem).data('callback');
            if(callback && $.isFunction(callback)) return callback.apply(this,arguments);
            window.setTimeout(function(){
                if(!nextURL || nextURL.indexOf('/sign_')>=0) nextURL='/index';
                window.location.href=nextURL;
            },2000);
        },
        error: function(xhr){cc.close()}
    };
    $.ajax(ajaxOptions);
}
/**
 * ajax注册
 * @param {string|object} elem 
 * @param {string} nextURL 
 */
function signUp(elem,nextURL){
    var cc=checkSubmitBtnWithCloser(elem);
    if(cc.submited)return;
    if(nextURL==null) nextURL=window.location.href;
    var data=$(elem).serializeArray();
    var url=$(elem).attr('action')||BASE_URL+'/sign_up';
    var ajaxOptions={
        url: url,
        type: 'POST',
        data: data,
        dataType: 'json',
        close: cc.close,
        success: function(r){
            if(r.Code!=1) {
                if(App.captchaHasError(r.Code) && r.Data && typeof(r.Data.captchaName)!='undefined' && r.Data.captchaName && $(elem).find('input[name="'+r.Data.captchaName+'"]').length<1){ 
                    return captchaDialog(r,ajaxOptions);
                }
                cc.close();
                showMsg({text:r.Info,type:'error'});
                return renewCaptcha(elem,r);
            }
            cc.close();
            showMsg({text:r.Info,type:'success'});
            var callback=$(elem).data('callback');
            if(callback && $.isFunction(callback)) return callback.apply(this,arguments);
            window.setTimeout(function(){
                if(!nextURL || nextURL.indexOf('/sign_')>=0) nextURL='/index';
                window.location.href=nextURL;
            },2000);
        },
        error: function(xhr){cc.close()}
    };
    $.ajax(ajaxOptions);
}
/**
 * 回车键操作
 */
function keyEnter(elem,callback){
    if($(elem).length<1) return;
    $(elem).on('keyup',function(event){
        if(event.keyCode==13) callback.apply(this);
    });
}
/**
 * 显示消息提示
 * @param {object} options 
 * @param {boolean} sticky 
 */
function showMsg(options,sticky){
    App.message(options,sticky);
}
/**
 * 关闭消息提示
 */
function closeMsg(){
    $.gritter.removeAll();
}
/**
 * 刷新验证码
 * @param {object} form 
 * @param {object} resp 
 */
function renewCaptcha(form,resp){
    if(form==null) return;
    App.captchaUpdate(form,resp);
};
/**
 * 显示图片验证码输入框
 * @param {object} resp ajax响应json对象
 * @param {object} ajaxOptions ajax选项
 * @author Wenhui Shen
 * @returns {boolean} 是否显示验证码输入框
 * @example
 * var ajaxOptions={
 *      url: window.location.href,
 *      type: 'POST',
 *      data: $('form').serializeArray(),
 *      dataType: 'json',
 *      success: function(resp){
 *          if(App.captchaHasError(resp.Code)){
 *              return captchaDialog(resp,ajaxOptions);
 *          }
 *          //Your code here...
 *      },
 *      error: function(xhr){
 *          //Your code here...
 *      }
 * };
 * $.ajax(ajaxOptions);
 */
function captchaDialog(resp,ajaxOptions){
    if(typeof(resp.Data)=='undefined'||typeof(resp.Data.captchaType)=='undefined'){
        showMsg({text:resp.Info,type:'error'});
        closeLoadingFunction(ajaxOptions);
        return false;
    }
    switch(resp.Data.captchaType){
        case 'api':return apiCaptchaDialog(resp,ajaxOptions);
        case 'go':return goCaptchaDialog(resp,ajaxOptions);
        default:return defaultCaptchaDialog(resp,ajaxOptions);
    }
}
function closeLoadingFunction(r){
    if(r&&typeof(r)=='object'&&typeof(r.close)=='function') r.close();
}
function goCaptchaDialog(resp,ajaxOptions){
    if(typeof(resp.Data)=='undefined' || typeof(resp.Data.driver)=='undefined'){
        showMsg({text:resp.Info,type:'error'});
        closeLoadingFunction(ajaxOptions);
        return false;
    }
    if(typeof(ajaxOptions.postByCaptchaDialog)!='undefined' && ajaxOptions.postByCaptchaDialog){
        showMsg({text:resp.Info,type:'error'});
    }
    var captchaName=resp.Data.captchaName,
        htmlCode=resp.Data.html
        jsInit=resp.Data.jsInit;
    for(var i=0;i<resp.Data.cssURLs.length;i++){
        var cssURL=resp.Data.cssURLs[i];
        if(cssURL&&$('link[href*="'+cssURL+'"]').length<1){
            $('head').append('<link rel="stylesheet" href="'+ASSETS_URL+cssURL+'?t='+BUILD_TIME+'">');
        }
    }
    for(var i=0;i<resp.Data.jsURLs.length;i++){
        var jsURL=resp.Data.jsURLs[i];
        if(jsURL&&$('script[src*="'+jsURL+'"]').length<1){
            $('body').append('<script src="'+ASSETS_URL+jsURL+'?t='+BUILD_TIME+'" type="text/javascript"></script>');
        }
    }
    var formHTML='<form method="post" id="dialog-retry-captcha">\
        '+htmlCode+'\
    </form>';
    var done=function(dialogRef) {
        var $form = $('#dialog-retry-captcha');
        var vcode = $form.find('[name="'+captchaName+'"]').val();
        postCaptchaDialogData(resp, ajaxOptions, vcode, null, captchaName, null);
        dialogRef.close();
        closeLoadingFunction(ajaxOptions);
    }
    App.dialog().show({
        title: App.t('行为验证'),
        type:'type-primary',
        size:'size-small',
        message:formHTML,
        nl2br:false,
        // closeByBackdrop:false,
        onshown: function(d){
            d.$modalDialog.css({"max-width":"340px"});
            d.$modalBody.css({"padding":0});
            var $input = $('#dialog-retry-captcha').find('[name="'+captchaName+'"]');
            if(!$input.hasClass('wg-cap-static')) $input.addClass('wg-cap-static');
            if(jsInit) {
                eval('var initCaptchaGo='+jsInit);
                initCaptchaGo(function(){
                    App.message({text:App.t('验证成功'),type:'success'});
                    done(d);
                });
                d.$modalBody.find('.wg-cap-wrap').css({"border":0,"border-top-left-radius":0,"border-top-right-radius":0});
            }
        },onhide: function(d){
            closeLoadingFunction(ajaxOptions);
        }/*,
        buttons: [{
            id: 'captchaDialogBtnSubmit',
            label: App.t('提交'),
            icon: 'fa fa-check',
            cssClass: 'btn-primary mg-r-10',
            action: done},{
            label: App.t('取消'),
            icon: 'fa fa-times',
            action: function(dialogRef) {
                dialogRef.close();
            }
        }]*/
    });
    return true;
}
function apiCaptchaDialog(resp,ajaxOptions){
    if(typeof(resp.Data)=='undefined' || typeof(resp.Data.provider)=='undefined'){
        showMsg({text:resp.Info,type:'error'});
        closeLoadingFunction(ajaxOptions);
        return false;
    }
    if(typeof(ajaxOptions.postByCaptchaDialog)!='undefined' && ajaxOptions.postByCaptchaDialog){
        showMsg({text:resp.Info,type:'error'});
    }
    var jsURL=resp.Data.jsURL,
        captchaName=resp.Data.captchaName,captchaIdent=resp.Data.captchaIdent,
        htmlCode=resp.Data.html,jsInit=resp.Data.jsInit;
    if(jsURL&&$('script[src="'+jsURL+'"]').length<1){
        $('body').append('<script src="'+jsURL+'" type="text/javascript"></script>');
    }
    var formHTML='<form method="post" id="dialog-retry-captcha">\
        <div class="form-group mg-b-0" style="position:relative">\
            <div class="captcha-loading text-center" id="'+resp.Data.locationID+'-loading" style="position:relative"><i class="fa fa-spinner fa-spin"></i> '+App.t('验证加载中，请稍候...')+'</div>\
            <div style="position:relative;z-index:2;margin-top:-20px;min-height:20px">'+htmlCode+'</div>\
        </div><!-- form-group -->\
    </form>';
    var done=function(dialogRef) {
        var $form = $('#dialog-retry-captcha');
        var vcode = $form.find('[name="'+captchaName+'"]').val();
        var idVal = $form.find('[name="'+captchaIdent+'"]').val();
        postCaptchaDialogData(resp, ajaxOptions, vcode, idVal, captchaName, captchaIdent);
        dialogRef.close();
        closeLoadingFunction(ajaxOptions);
    }
    /*
    var captchaID=resp.Data.captchaID,jsCallback=resp.Data.jsCallback;
    if(jsCallback) {
        eval('window["apiCaptchaCallback'+captchaID+'"]='+jsCallback);
        var oldDone=done;
        done=function(d){
            window["apiCaptchaCallback"+captchaID](function(){
                oldDone(d)
            });
        };
    }
    */
    App.dialog().show({
        title:App.t('人机验证'),
        type:'type-primary',
        size:'size-small',
        message:formHTML,
        nl2br:false,
        closeByBackdrop:false,
        onshown: function(d){
            if(jsInit) eval(jsInit);
            $('#dialog-retry-captcha').on('submit',function(e){
                e.preventDefault();
                d.getButton('captchaDialogBtnSubmit').trigger('click');
            });
        },onhide: function(d){
            closeLoadingFunction(ajaxOptions);
        },
        buttons: [{
            id: 'captchaDialogBtnSubmit',
            label: App.t('提交'),
            icon: 'fa fa-check',
            cssClass: 'btn-primary mg-r-10',
            action: done},{
            label: App.t('取消'),
            icon: 'fa fa-times',
            action: function(dialogRef) {
                dialogRef.close();
                closeLoadingFunction(ajaxOptions);
            }
        }]
    });
    return true;
}
function defaultCaptchaDialog(resp,ajaxOptions){
    if(typeof(resp.Data)=='undefined' || typeof(resp.Data.captchaURL)=='undefined'){
        showMsg({text:resp.Info,type:'error'});
        closeLoadingFunction(ajaxOptions);
        return false;
    }
    if(typeof(ajaxOptions.postByCaptchaDialog)!='undefined' && ajaxOptions.postByCaptchaDialog){
        showMsg({text:resp.Info,type:'error'});
    }
    var captchaName=resp.Data.captchaName,captchaURL=resp.Data.captchaURL;
    var captchaID=resp.Data.captchaID,captchaIdent=resp.Data.captchaIdent;
    var formHTML='<form method="post" id="dialog-retry-captcha">\
        <div class="form-group mg-b-0">\
            <div class="input-group">\
              <input type="text" placeholder="'+App.t('图像验证码')+'" value="" name="'+captchaName+'" required="required" data-toggle="tooltip" title="'+App.t('请输入验证码')+'" class="form-control" id="dialog-captcha-code">\
              <span class="input-group-addon pd-0"><img id="captchaImage" src="'+captchaURL+'" alt="Captcha image" onclick="this.src=this.src.split(\'?\')[0]+\'?reload=\'+Math.random();" onerror="if(this.src.indexOf(\'?reload=\')!=-1 && confirm(\''+App.t('页面验证码已经失效，必须重新请求当前页面。确定要刷新本页面吗？')+'\')) window.location.reload();" style="cursor:pointer"><input type="hidden" id="dialog-captcha-id" name="'+captchaIdent+'" value="'+captchaID+'"></span>\
            </div>\
            <small id="dialog-captcha-code-help" class="form-text text-muted">'+App.t('请输入图片上显示的验证码')+'</small>\
        </div><!-- form-group -->\
    </form>';
    App.dialog().show({
        title:App.t('操作验证'),
        type:'type-primary',
        size:'size-small',
        message:formHTML,
        nl2br:false,
        closeByBackdrop:false,
        onshown: function(d){
            var ipt=$('#dialog-captcha-code');
            ipt.focus();
            ipt.select();
            $('#dialog-retry-captcha').on('submit',function(e){
                e.preventDefault();
                d.getButton('captchaDialogBtnSubmit').trigger('click');
            });
        },onhide: function(d){
            closeLoadingFunction(ajaxOptions);
        },
        buttons: [{
            id: 'captchaDialogBtnSubmit',
            label: App.t('提交'),
            icon: 'fa fa-check',
            cssClass: 'btn-primary mg-r-10',
            action: function(dialogRef) {
                var vcode = $('#dialog-captcha-code').val();
                var idVal = $('#dialog-captcha-id').val();
                postCaptchaDialogData(resp, ajaxOptions, vcode, idVal, captchaName, captchaIdent);
                dialogRef.close();
                closeLoadingFunction(ajaxOptions);
            }
        },{
            label: App.t('取消'),
            icon: 'fa fa-times',
            action: function(dialogRef) {
                dialogRef.close();
                closeLoadingFunction(ajaxOptions);
            }
        }]
    });
    return true;
}
function postCaptchaDialogData(resp, ajaxOptions, vcode, idVal, captchaName, captchaIdent){
    var isAjaxForm = ('ajaxFormObject' in ajaxOptions) && ajaxOptions.ajaxFormObject;
    if(!isAjaxForm || (resp.Zone && $(ajaxOptions.ajaxFormObject).find('[name="'+resp.Zone+'"]').length<1)){
        if(!ajaxOptions.data) ajaxOptions.data={};
        switch(typeof(ajaxOptions.data)){
            case 'string':
            ajaxOptions.data+='&'+captchaName+'='+encodeURIComponent(vcode);
            if(captchaIdent) ajaxOptions.data+='&'+captchaIdent+'='+encodeURIComponent(idVal);
            break;
            default:
            if($.isArray(ajaxOptions.data)){
                var existsCode=false,existsIdent=false;
                var codeInfo={name:captchaName,value:vcode};
                var identInfo=captchaIdent?{name:captchaIdent,value:idVal}:null;
                for(var i=0; i<ajaxOptions.data.length; i++){
                    if(ajaxOptions.data[i].name==captchaName){
                        ajaxOptions.data[i]=codeInfo;
                        existsCode=true;
                    }
                    if(captchaIdent){
                        if(ajaxOptions.data[i].name==captchaIdent){
                            ajaxOptions.data[i]=identInfo;
                            existsIdent=true;
                        }
                        if(existsIdent && existsCode) break;
                    }else{
                        if(existsCode) break;
                    }
                }
                if(!existsIdent&&identInfo) ajaxOptions.data.push(identInfo);
                if(!existsCode) ajaxOptions.data.push(codeInfo);
                break;
            }
            ajaxOptions.data[captchaName] = vcode;
            if(captchaIdent)ajaxOptions.data[captchaIdent] = idVal;
        }
    }
    ajaxOptions.postByCaptchaDialog=true;
    if(isAjaxForm){
        var form=$(ajaxOptions.ajaxFormObject);
        var codeIpt=form.find('[name="'+captchaName+'"]');
        if(codeIpt.length>0){
            codeIpt.val(vcode);
            if(captchaIdent) form.find('[name="'+captchaIdent+'"]').val(idVal);
        }
        form.trigger('submit');
    }else{
        $.ajax(ajaxOptions);
    }
}
// 登录
function signInDialog(callback,r){
    if($('#modal-sign-in').length>0) return $('#modal-sign-in').modal('show');
    $.get('/sign_in',{modal:1,next:window.location.href},function(r){
        $('body').append(r);
        closeLoadingFunction(r);
        if(callback!=null && $.isFunction(callback)){
            $('#modal-sign-in-form').data('callback',function(){
                $('#modal-sign-in').modal('hide');
                return callback.apply(this,arguments);
            });
        }
        $('#modal-sign-in').modal('show');
        $('#modal-sign-in').on('shown.bs.modal',function(){
            $('#modal-sign-in-form').find('input[name="name"]').focus();
        });
    },'html').fail(function(){
        closeLoadingFunction(r);
    });
}
/**
 * 输入验证码
 * @param {string} objectName 物件名(短信/邮件)
 * @param {object} formData 表单数据对象
 * @param {cloure} onSuccess 提交成功后的回调
 */
function inputVerifyCode(objectName,formData,postURL,onSuccess){
  if(postURL==null) postURL=window.location.href;
  var formHTML='<form method="post" id="dialog-vcode-form">\
  <div class="form-group mg-b-0">\
    <input type="text" placeholder="'+App.t('请输入%s验证码', objectName)+'" value="" name="code" required="required" data-toggle="tooltip" title="'+App.t('请输入%s验证码', objectName)+'" class="form-control" id="dialog-input-vcode">\
    <small id="dialog-input-vcode-help" class="form-text text-muted">'+App.t('请输入您接收到的%s验证码', objectName)+'</small>\
  </div>\
</form>';
    App.dialog().show({
        title:App.t('请输入%s验证码',objectName),
        type:'type-primary',
        size:'size-small',
        message:formHTML,
        nl2br:false,
        closeByBackdrop:false,
        onshown:function(d){
          var ipt=$('#dialog-input-vcode');
          ipt.focus();
          ipt.select();
          $('#dialog-vcode-form').on('submit',function(e){
              e.preventDefault();
              d.getButton('inputVerifyCodeBtnSubmit').trigger('click');
          });
        },
        buttons: [{
            id: 'inputVerifyCodeBtnSubmit',
            label: App.t('提交'),
            icon: 'fa fa-check',
            cssClass: 'btn-primary mg-r-10',
            action: function(dialogRef) {
                var vcode = $('#dialog-input-vcode').val();
                var data;
                var pushArray=function(){
                  data.push({name:"vcode",value:vcode});
                  data.push({name:"verify",value:1});
                };
                switch(typeof(formData)){
                  case 'object':
                  if(!$.isArray(formData)){
                    data = {};
                    for(var i in formData){
                        data[i]=formData[i];
                    }
                    data["vcode"]=vcode;
                    data["verify"]=1;
                    break;
                  }

                  data = [];
                  for(var i=0;i<formData.length;i++){
                      data.push(formData[i]);
                  }
                  pushArray();
                  break;

                  case 'string':
                  data=$(formData).serializeArray();
                  pushArray();
                  break;

                  default:
                  data=[];
                  pushArray();
                  break;
                }
                $.post(postURL,data,function(r){
                  if(r.Code==1){
                    dialogRef.close();
                    App.message({text:r.Info,type:'success'});
                    if(onSuccess!=null) return onSuccess.apply(this,arguments);
                    window.setTimeout(function(){window.location.reload()},2000);
                    return;
                  }
                  App.message({text:r.Info,type:'error'});
                },'json');
            }
        },{
            label: App.t('取消'),
            icon: 'fa fa-times',
            action: function(dialogRef) {
                dialogRef.close();
            }
        }]
    });
}
/**
 * 身份验证对话框
 * @param {string} authType 认证类型(password/email/mobile)
 * @param {string} objectName 发送物件名称(邮件/短信)
 * @param {string} formElem 表单元素(#id)
 * @param {number} waitingSeconds 每次发送验证码的等待时间
 * @param {string} postURL 修改密码的提交网址
 */
function authDialog(authType,objectName,formElem,postURL){
if(postURL==null) postURL=window.location.href;
$.get(BASE_URL+'/user/profile/_authentication',{type:authType},function(r){
    var html=r.Data.html;
    var btnLabel=App.t('验证密码'),btnIcon='check';
    var waitingSeconds=r.Data.waitingSeconds;
    switch(r.Data.type){
        case 'email':btnLabel=App.t('发送邮件');btnIcon='paper-plane';break;
        case 'mobile':btnLabel=App.t('发送短信');btnIcon='paper-plane';break;
    }
    App.dialog().show({
        id:'authentication-dialog',
        title:App.t('身份验证'),
        type:'type-primary',
        message:html,
        nl2br:false,
        closeByBackdrop:false,
        onshown: function(dialogRef){
          var $form=$('#authentication-form');
          var ipt=$form.find('input:not(readonly):first');
          ipt.focus();
          ipt.select();
          $form.prev('ul.nav').find('a.nav-link').on('click',function(){
            authDialog($(this).data('type'),$(this).data('object-name'),formElem,postURL);
            dialogRef.close();
          });
        },
        buttons: [{
            id: 'authentication-btn-submit',
            label: btnLabel,
            icon: 'fa fa-'+btnIcon,
            cssClass: 'btn-primary mg-r-10',
            action: function(dialogRef) {
                var btn = this;
                var formData=$(formElem).serializeArray();
                var data=$('#authentication-form').serializeArray();
                for(var i=0;i<data.length;i++) formData.push(data[i]);
                var sep=postURL.indexOf('?')==-1?'?':'&';
                //验证旧密码时没有发送信件的步骤
                var step=(authType=='password')?'verify':'send';
                if(step=='send')btnDisabled(btn,'<i class="fa fa-spinner fa-spin"></i> '+App.t('发送中，请稍候...'));
                $.post(postURL+sep+'step='+step,formData,function(r){
                    if(step=='send')btnEnabled(btn);
                    renewCaptcha('#authentication-form',r);
                    if(r.Code!=1) return App.message({text:r.Info,type:'error'});
                    if(step=='send')waiting(btn,waitingSeconds);
                    if(r.Data.nextStep=='verify'){
                        App.message({text:r.Info+','+App.t('请输入您收到的验证码'),type:'success'});
                        return inputVerifyCode(objectName,formData,postURL+sep+'step=verify',function(){
                            dialogRef.close();
                        });
                    }
                    App.message({text:r.Info,type:'success'});
                    dialogRef.close();
                },'json');
            }
        },{
            label: App.t('取消'),
            icon: 'fa fa-times',
            action: function(dialogRef) {
                dialogRef.close();
            }
        }]
    });
},'json');
}
/**
 * 获取评论列表
 * @param {integer} page 
 * @param {string|object} boxElem 
 * @param {string} url 
 */
function commentList(page,boxElem,url,isReplyList){
    if(page==null) page=1;
    if(boxElem==null) boxElem='#comment-list-container';
    if(isReplyList==null) isReplyList=false;
    var box=$(boxElem);
    if(!box||box.length<1) return;
    if(url==null) url=box.data('url');
    if(!url) return;
    loadCommentList(box,url,{page:page,html:1,_pjax:'true'},isReplyList,$('#headPanel').height()+90);
}

/**
 * 加载评论列表
 * @param {object} box 
 */
 function loadCommentList(box,url,data,isReplyList,offsetY) {
    $.get(url,data||{},function(r){
        if(r.Code!=1) return showMsg(r.Info);
        box.html(r.Data.html);
        bindCommentList(box,isReplyList,offsetY);
        if(!isReplyList && r.Data.pagination){
            $('.comment-count-num').text(r.Data.pagination.rows);
        }
    },'json');
}
/**
 * 绑定评论列表事件
 * @param {object} box 
 * @param {boolean} isReplyList
 * @param {float} offsetY 
 */
function bindCommentList(box,isReplyList,offsetY) {
    if(!box || box.length<1) return;
    if(offsetY===null) offsetY=0;
    if(isReplyList===null) isReplyList=false;
    var replaceState = box.data('replace-state')||false;
    box.find('[data-reply-id]').on('click',function(){
        var replyId=$(this).data('reply-id');
        var replyBox=$('#reply-list-box-'+replyId);
        commentList(1,replyBox[0],replyBox.data('url'),true);
    });
    attachContype(box);
    box.find('[data-comment-like-id]').on('click',function(){
        var id=$(this).data('comment-like-id'),me=$(this);
        $.post(BASE_URL+'/article/comment_like',{id:id},function(r){
            if(r.Code<1){
                App.message({text:r.Info,type:"error"});
                return;
            }
            me.removeClass('tx-gray-light').addClass('tx-success');
            var numElem=me.find('.comment-likes');
            if(numElem.length>0){
                numElem.text(numElem.text()*1+1);
            }
        },'json');
    });
    box.find('.pagination a[data-pjax="true"]').on('click',function(e){
        e.preventDefault();
        if(replaceState) App.replaceState(null,'',$(this).attr('href'));
        loadCommentList(box,$(this).attr('href'),null,isReplyList,offsetY);
        $('html, body').animate({scrollTop: box.offset().top-offsetY}, 0);
    });
}
/**
 * 初始化评论表单事件
 * @param {string|object} formElem 
 */
function initCommentForm(formElem){
    if(formElem==null) formElem='#comment-post-form';
    if($(formElem).length<1) return;
    if($(formElem).data('attached')) return;
    $(formElem).data('attached',true);
    $(formElem).find('[required]').on('keyup blur',function(){
        var v=$.trim($(this).val());
        if(v==''){
            $(this).addClass('form-control-warning').parent().addClass('has-warning');
        }else{
            $(this).removeClass('form-control-warning').parent().removeClass('has-warning')
        }
    });
    var f=function(){
        var invalidInput=$(formElem).find('.form-control-warning[required]');
        if(invalidInput.length>0) {
            invalidInput.focus();
            showMsg(invalidInput.attr('placeholder'));
            return;
        }
        var url=$(formElem).attr('action');
        var valid=true;
        $(formElem).find('[required]').each(function(){
            var v=$.trim($(this).val());
            if(v==''){
                if(valid){
                    $(this).focus();
                    showMsg($(this).attr('placeholder'));
                }
                valid=false;
                $(this).addClass('form-control-warning').parent().addClass('has-warning');
            }
        });
        if(!valid) return;
        $.post(url,$(formElem).serializeArray(),function(r){
            renewCaptcha($(formElem),r);
            showMsg(r.Info);
            if(r.Code>0){
                if($('.show-after-comment.mission-uncompleted').length>0){
                    window.setTimeout(function(){
                        window.location.reload();
                    },2000);
                    return;
                }
                if(/^modal-/.test($(formElem).attr('id'))){
                    $(formElem).parents('.modal').modal('hide');
                }
                var replyIdE=$(formElem).find('#reply-id');
                var replyId=replyIdE.data('root-id')||0;
                if(replyId && $('#reply-list-box-'+replyId).length>0){
                    var replyBox = $('#reply-list-box-'+replyId);
                    commentList(1,replyBox[0],replyBox.data('url'),true);
                }else{
                    commentList(1);
                }
                $(formElem).find('[name=content]').val('');
                $(formElem).find('[name=code]').val('');
            }
        },'json');
    };
    $(formElem).on('submit', function(event){
        event.preventDefault();
        f.apply(this,arguments);
    });
    $(formElem).find('[data-form]:not(:submit)').on('click',f);
}
function initPageAllCommentForm(){
    initCommentForm();
    initCommentForm('#modal-comment-post-form');
    $('#modal-comment-form').on('show.bs.modal', function (event) {
        var button = $(event.relatedTarget);
        var commentId = button.data('comment-id')||'0';
        var rootId = button.data('root-id')||'0';
        var modal = $(this);
        var replyIdE = modal.find('#reply-id');
        replyIdE.val(commentId);
        replyIdE.data('root-id',rootId);
        var blk=modal.find('#reply-comment-blockquote');
        if(commentId>0){//回复评论
            blk.show();
            blk.find('#reply-comment-content').html($(button.data('comment-content')).html());
            blk.find('#reply-comment-author').html(button.data('comment-author'));
        }else{
            blk.hide();
            blk.find('#reply-comment-content').html('');
            blk.find('#reply-comment-author').html('');
        }
    });
}
/**
 * 关注用户
 * @param {string|object} a 
 */
function follow(a){
    var uid=$(a).data('following-uid');
    var params={uid:uid};
    var unfollow=$(a).hasClass('btn-outline-primary');
    if(!unfollow) params['unfollow']=1;
    $.post(BASE_URL+'/user/profile/following',params,function(r){
        if(r.Code<1){
            showMsg({text:r.Info,type:'error'});
            if(r.Code==App.status.NotLoggedIn) window.location.href=r.URL+'?next='.encodeURLComponent(window.location.href);
            return;
        }
        showMsg({text:r.Info,type:'success'});
        if(unfollow){
            $(a).removeClass('btn-outline-primary').addClass('btn-outline-gray-lighter active');
            $(a).text($(a).data('i18n-followed'));
        }else{
            $(a).removeClass('btn-outline-danger').removeClass('btn-outline-gray-lighter active').addClass('btn-outline-primary');
            $(a).text($(a).data('i18n-following'));
        }
    },'json');
}
/**
 * 设置按钮关注状态
 * @param {string|object} a 
 */
function setFollowStatus(a){
    var uid=$(a).data('following-uid');
    if(uid==CUSTOMER.ID){
        return $(a).replaceWith('<small>'+App.t('这是我')+'</small>');
    }
    $.post(BASE_URL+'/user/profile/is_followed',{uid:uid},function(r){
        if(r.Code>0){
            $(a).removeClass('btn-outline-primary').addClass('btn-outline-danger');
            $(a).text($(a).data('i18n-unfollow'));
        }
    },'json');
}
function btnDisabled(btn,innerHTML){
    var ori=btn.data('original-innerhtml');
    if(!ori){
      ori=btn.html();
      btn.data('original-innerhtml',ori);
    }
    btn.prop('disabled',true);
    btn.html(innerHTML);
}
function btnEnabled(btn){
    var ori=btn.data('original-innerhtml');
    btn.html(ori);
    btn.prop('disabled',false);
}
function waiting(btn,waitingSeconds){
  if(btn.prop('disabled')) return;
  if(waitingSeconds==null) waitingSeconds=60;
  //var btn=$('#binding-form').find('button:submit');
  var labelId='timedown-'+String(Math.random()).replace('.','');
  btnDisabled(btn,'<i class="fa fa-spinner fa-spin"></i> <span id="'+labelId+'">'+waitingSeconds+App.t('秒后可重试')+'</span>')
  var timedown=$('#'+labelId);
  var n=waitingSeconds;
  var t=window.setInterval(function(){
    n--;
    if(n<=0){
      window.clearInterval(t);
      btnEnabled(btn);
      return;
    }
    timedown.text(n+App.t('秒后可重试'));
  },1000);
}
function countWords(a,countElem){
    if(countElem==null){
        var id=$(a).attr('id'),name=$(a).attr('name');
        if($('#'+id+'WordsCount').length>0) countElem='#'+id+'WordsCount';
        else if($('#'+name+'WordsCount').length>0) countElem='#'+name+'WordsCount';
        else throw 'Variable countElem is undefined';
    }
    $(countElem).text(Number($(a).attr('maxlength'))-$(a).val().length);
}

function checkSubmitBtnWithCloser(a){
    var $submit = $(a).find(':submit'),close=function(){};
    if($submit.length>0){
        if($submit.prop('disabled')) return {submited:true,close:close};
        var $icon=$submit.children('.fa');
        $submit.prop('disabled',true);
        if($icon.length>0)$icon.addClass('fa-refresh fa-spin');
        close=function(){
            $submit.prop('disabled',false);
            if($icon.length>0)$icon.removeClass('fa-refresh fa-spin');
        };
    }
    return {submited:false,close:close};
}
function ajaxForm(a,onSuccess,onFailure){
    var cc=checkSubmitBtnWithCloser(a);
    if(cc.submited)return;
    var opts={
        ajaxFormObject: a,
        type: String($(a).attr('method')).toLowerCase()=='post'?'post':'get',
        dataType: 'json',
        data: {},
        url: $(a).attr('action'),
        close: cc.close,
        success: function(r){
            onAjaxRespond(a,r,opts,onSuccess,onFailure);
        },error: function(xhr){
            cc.close();
        }
    };
    $(a).ajaxForm(opts);
}
function onAjaxRespond(form,r,ajaxOptions,onSuccess,onFailure){
    if(r.Code==1){
        renewCaptcha(form,r);
        if(onSuccess!=null&&$.isFunction(onSuccess)) onSuccess(r);
        closeLoadingFunction(ajaxOptions);
        return showMsg({text:r.Info,type:'success'});
    }
    if(r.Code==App.status.NotLoggedIn){
        return signInDialog(function(){},ajaxOptions);
    }
    if(App.captchaHasError(r.Code)) return captchaDialog(r,ajaxOptions);
    showMsg({text:r.Info,type:'error'});
    renewCaptcha(form,r);
    if(onFailure!=null&&$.isFunction(onFailure)) onFailure(r);
    closeLoadingFunction(ajaxOptions);
}
function hasScroll(el, direction) {
    var eleScroll = (!direction || direction === 'vertical') ? 'scrollTop' : 'scrollLeft';
    var result = !!el[eleScroll];
    if (!result) {
        el[eleScroll] = 1;
        result = !!el[eleScroll];
        el[eleScroll] = 0;
    }
    return result;
}
/**
 * 投诉
 * @param {string|object} elem
 */
function complaint(elem,close){
    var data=$(elem).serializeArray();
    var url=$(elem).attr('action')||BASE_URL+'/complaint';
    var ajaxOptions={
        url: url,
        type: 'POST',
        data: data,
        dataType: 'json',
        close: close,
        success: function(r){
            if(r.Code!=1) {
                if(App.captchaHasError(r.Code) && $(elem).find('input[name="'+r.Zone+'"]').length<1){ 
                    return captchaDialog(r,ajaxOptions);
                }
                close();
                showMsg({text:r.Info,type:'error'});
                return renewCaptcha(elem,r);
            }
            close();
            showMsg({text:r.Info,type:'success'});
            var callback=$(elem).data('callback');
            if(callback && $.isFunction(callback)) return callback.apply(this,arguments);
        },
        error: function(xhr){
            close();
        }
    };
    $.ajax(ajaxOptions);
}
/**
 * 投诉对话框
 * @param {string} targetElem 触发元素
 * @param {cloure} onSuccess 提交成功后的回调
 */
function complaintModal(targetElem,onSuccess){
  var postURL=$(targetElem).data('complaint-url');
  if(!postURL) return;
  var objectName=$(targetElem).data('complaint-name');
  var typeList=$(targetElem).data('complaint-types')||[];
  var typeHTML='';
  for(var i=0;i<typeList.length;i++){
    var type=typeList[i],seled=i==0?' checked="checked"':'';
    typeHTML+='<label class="custom-control custom-radio">\
        <input type="radio" class="custom-control-input" value="'+type.K+'" name="type"'+seled+'>\
        <span class="custom-control-indicator"></span>\
        <span class="custom-control-description">'+type.V+'</span>\
    </label>';
  }
  var formHTML='<form method="post" id="dialog-complaint-form" action="'+postURL+'">\
  <div class="form-group mg-b-0">'+typeHTML+'</div>\
  <div class="form-group mg-b-0">\
    <textarea placeholder="'+App.t('投诉内容')+'" value="" name="content" required="required" data-toggle="tooltip" title="'+App.t('请输入投诉内容')+'" class="form-control" id="dialog-input-complaint-content"></textarea>\
    <small id="dialog-input-complaint-content-help" class="form-text text-muted">'+App.t('请输入投诉内容')+'</small>\
  </div>\
</form>';
    App.dialog().show({
        title:App.t('投诉《%s》',objectName),
        type:'type-primary',
        size:'size-middle',
        message:formHTML,
        nl2br:false,
        closeByBackdrop:true,
        onshown:function(d){
          var ipt=$('#dialog-input-complaint-content');
          ipt.focus();
          ipt.select();
          $('#dialog-complaint-form').on('submit',function(e){
              e.preventDefault();
              d.getButton('complaintFormBtnSubmit').trigger('click');
          });
        },
        buttons: [{
            id: 'complaintFormBtnSubmit',
            label: App.t('提交'),
            icon: 'fa fa-check',
            cssClass: 'btn-primary mg-r-10',
            action: function(dialogRef) {
                var $submit = $('#complaintFormBtnSubmit'),$icon=$submit.children('.fa');
                if($submit.prop('disabled'))return;
                var ipt=$('#dialog-input-complaint-content');
                var val=$.trim(ipt.val());
                if(val==''){
                    showMsg({text:App.t('请输入投诉内容'),type:'warning'});
                    return;
                }
                ipt.val(val);
                $('#dialog-complaint-form').data('callback',function(r){
                    if(r.Code==1){
                      dialogRef.close();
                      if(onSuccess!=null) return onSuccess.apply(this,arguments);
                    }
                });
                $submit.prop('disabled',true);
                $icon.addClass('fa-refresh fa-spin');
                var close=function(){
                    $submit.prop('disabled',false);
                    $icon.removeClass('fa-refresh fa-spin');
                };
                complaint('#dialog-complaint-form',close);
            }
        },{
            label: App.t('取消'),
            icon: 'fa fa-times',
            action: function(dialogRef) {
                dialogRef.close();
            }
        }]
    });
}
/**
 * 获取用于恢复密码的验证码
 * @param {string} targetElem 触发元素
 * @param {cloure} onSuccess 提交成功后的回调
 */
function forgotModal(targetElem,onSuccess){
    var postURL=$(targetElem).data('submit-url');
    var formTmpl=$(targetElem).data('form-tmpl');
    var nameSync=$(targetElem).data('name-sync');
    if(!postURL) return;
    var formId = 'dialog-recv-recover-code-form', elem = '#'+formId;
    var formHTML='<form method="post" id="'+formId+'" action="'+postURL+'">';
    formHTML+=$(formTmpl).html();
    formHTML+='</form>';
    App.dialog().show({
        title:App.t('找回密码：获取验证码'),
        type:'type-primary',
        size:'size-middle',
        message:formHTML,
        nl2br:false,
        closeByBackdrop:true,
        onshown:function(d){
            var $f = $(elem);
            if(nameSync){
                var v=$(nameSync).val();
                if(v) $f.find('input[name="name"]').val(v);
            }
            var ipt=$f.find('input[type=text]:first');
            ipt.focus();
            ipt.select();
            $f.on('submit',function(e){
                e.preventDefault();
                d.getButton('recvRecoverCodeFormBtn').trigger('click');
            });
            $f.find('select[name=type]').off().on('change',function(){
                switch($(this).val()){
                    case 'mobile':
                        $f.find('input[name=account]').attr('placeholder',App.t('手机号码'));
                        break;
                    default:
                        $f.find('input[name=account]').attr('placeholder',App.t('E-mail地址'));
                        break;
                }
            });
        },
        buttons: [{
            id: 'recvRecoverCodeFormBtn',
            label: App.t('发送'),
            icon: 'fa fa-check',
            cssClass: 'btn-primary mg-r-10',
            action: function(dialogRef) {
                var $submit = $('#recvRecoverCodeFormBtn'),$icon=$submit.children('.fa');
                if($submit.prop('disabled'))return;
                $submit.prop('disabled',true);
                $icon.addClass('fa-refresh fa-spin');
                var close=function(){
                    $submit.prop('disabled',false);
                    $icon.removeClass('fa-refresh fa-spin');
                };
                var names = ['name', 'account'];
                for (var i = 0; i<names.length; i++) {
                    var name = names[i];
                    var ipt=$(elem).find('input[name="'+name+'"]');
                    var val=$.trim(ipt.val());
                    if(val==''){
                        var msg = '';
                        switch(name){
                            case 'name':
                                msg = App.t('请输入用户名');
                                break;
                            case 'account':
                                msg = App.t('请输入%s', ipt.attr('placeholder'));
                                break;
                            default:
                                msg = name;
                        }
                        showMsg({text:msg,type:'warning'});
                        ipt.focus();
                        close();
                        return;
                    }
                    ipt.val(val);
                }
                var ajaxOptions={
                    url: postURL,
                    type: 'POST',
                    data: $(elem).serializeArray(),
                    dataType: 'json',
                    close: close,
                    success: function(r){
                        if(r.Code!=1) {
                            if(App.captchaHasError(r.Code) && $(elem).find('input[name="'+r.Zone+'"]').length<1){
                                return captchaDialog(r,ajaxOptions);
                            }
                            close();
                            showMsg({text:r.Info,type:'error'});
                            return renewCaptcha(elem,r);
                        }
                        showMsg({text:r.Info,type:'success'});
                        if(nameSync){
                            var $f = $(elem);
                            $(nameSync).val($f.find('input[name="name"]').val());
                        }
                        if(onSuccess!=null) {
                            onSuccess.apply($(elem),arguments);
                        }
                        dialogRef.close();
                        close();
                    },
                    error: function(xhr){
                        close();
                    }
                };
                $.ajax(ajaxOptions);
            }
        },{
            label: App.t('取消'),
            icon: 'fa fa-times',
            action: function(dialogRef) {
                dialogRef.close();
            }
        }]
    });
}
function openForgotModal(a){
    forgotModal(a,function(r){
        $('#recvType').val(this.find('select[name=type]').val());
        $('#recvAccount').val(this.find('input[name=account]').val());
    });
}
/**
 * 重置密码
 * @param {string|object} elem
 * @param {string} nextURL
 */
function resetPassword(elem,nextURL){
    var data=$(elem).serializeArray();
    var url=$(elem).attr('action')||BASE_URL+'/forgot';
    var ajaxOptions={
        url: url,
        type: 'POST',
        data: data,
        dataType: 'json',
        success: function(r){
            if(r.Code!=1) {
                if(App.captchaHasError(r.Code) && $(elem).find('input[name="'+r.Zone+'"]').length<1){
                    return captchaDialog(r,ajaxOptions);
                }
                showMsg({text:r.Info,type:'error'});
                return renewCaptcha(elem,r);
            }
            showMsg({text:r.Info,type:'success'});
            var callback=$(elem).data('callback');
            if(callback && $.isFunction(callback)) return callback.apply(this,arguments);
            window.setTimeout(function(){
                if(!nextURL || nextURL.indexOf('/forgot')>=0) nextURL='/index';
                window.location.href=nextURL;
            },2000);
        },
        error: function(xhr){}
    };
    $.ajax(ajaxOptions);
}
function getColorByRandom(colorList) {
    var colorIndex = Math.floor(Math.random() * colorList.length);
    var color = colorList[colorIndex];
    colorList.splice(colorIndex, 1);
    return color;
}
function tagsRandomColor(elem, startIndex){
    if(!elem) elem=".tags-container > .tag-item";
    if(startIndex==null) startIndex=0;
    var colorList = ["orange","coral","darkorange","darksalmon","hotpink","#9dc6eb","#f8c471","#b9a3ef","#fdb1ca","#9dc6eb","#f8c471","#b9a3ef","#fdb1ca"];
    var tags = $(elem);
    var total = tags.length;
    for (var i = startIndex; i < total; i++) {
        var bgColor = getColorByRandom(colorList);
        tags.eq(i).css("background-color", bgColor);
    }
}
function backLastPage(index) {
    if(index==null){
        window.location.href=document.referrer;
        return;
    }
    window.history.go(index);
}
function setNavActive(){
    var act = $('#navbarMain .active');
    if(act.length<1)return;
    while(true){
        var className = act.attr('class');
        if(act.hasClass('dropdown-item')){
            act = act.closest('.dropdown-menu').prev('a:not(.active)');
        }else{
            act = act.closest('ul').prev('a:not(.active)');
        }
        if(act.length<1) return;
        if(act.hasClass('nav-link')) return act.addClass('active');
        act.addClass(className);
    }
}
function attachContype(container){
    if(typeof App.editor === 'undefined') return;
    App.editor.attachContype(container);
}
function commonInit($,App){
    if(window.errorMSG) App.message({title: App.i18n.SYS_INFO, text: App.ifTextNl2br(window.errorMSG), class_name: "danger"});
	if(window.successMSG) App.message({title: App.i18n.SYS_INFO, text: App.ifTextNl2br(window.successMSG), class_name: "success"});
    $(function(){
        $('[data-toggle="tooltip"]:not(data-original-title)').tooltip();
        $('[data-following-uid]').each(function(){
            setFollowStatus(this);
        });
        $('textarea[count-words]').on('keyup',function(){
            countWords(this,$(this).attr('count-words'));
        }).trigger('keyup');
        $('form[ajax-submit]').each(function(){
            ajaxForm(this);
        });
        App.fixedFooter('.footer-fix');
        App.bottomFloat('.auto-bottom-float',0,true);
        App.showRequriedInputStar();
        setNavActive();
        attachContype();
    })
}
if (typeof(jQuery) !== 'undefined' && typeof(App) !== 'undefined') {
    commonInit(jQuery,App);
}
