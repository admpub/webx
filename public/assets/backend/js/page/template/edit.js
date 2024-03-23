function openTemplateFile(obj,file){
    var $unsaved=$(obj).parent('li').siblings('li.active').find('.badge-unsaved');
    if($unsaved.length>0){
        var $fileObj=$unsaved.next('a');
        if(!confirm(App.t('文件“%s”尚未保存，此时切换到其它界面将会丢失已更改内容，确定要放弃未保存的内容吗？',$fileObj.data('file')))){
            return;
        }
        if($fileObj.data('new')){
            $fileObj.parent('li').remove();
        }else{
            $unsaved.remove();
        }
    }
    if(file==null) file=$(obj).data('file');
    var dir=$('#templateDirs').data('dir');
    var name=$('#templateDirs').data('name');
    $.get(window.location.pathname,{name:name,op:'getFileContent',dir:dir,file:file},function(r){
        if(r.Code!=1)return App.message({text:r.Info,type:'error'});
        var editor=$('#console').data('codemirror');
        $('#console').data('onceignorechange',true);
        editor.setValue(r.Data.content);
        editor.refresh();
        var $tgl = $('#main-container>.page-aside>.navbar-toggle')
        if($tgl.is(':visible')) $tgl.trigger('click');
        if(!$(obj).parent('li').hasClass('active')) $(obj).parent('li').addClass('active').siblings('li.active').removeClass('active');
        $('#currentOperateFile').html(file);
    },'json');
}
function saveTemplateFile(obj,file){
    var $fileObj = $('#templateDirs').find('li.active>a.editable');
    if(file==null) file=$fileObj.data('file');
    var isNew = $fileObj.data('new')||'';
    var dir=$('#templateDirs').data('dir');
    var name=$('#templateDirs').data('name');
    var save=function(confirmed){
        var data={
            name:name,
            op:'saveFileContent',
            dir:dir,
            file:file,
            content:$('#console').data('codemirror').getValue(),
            isNew:isNew
        };
        if(confirmed) data.confirmed = 1;
        $.post(window.location.pathname,data,function(r){
            if(r.Code!=1){
                if(r.Code==-109){ // 文件已经存在，确认是否覆盖
                    if(confirm(r.Info)){
                        return save(true);
                    }
                    return;
                }
                return App.message({text:r.Info,type:'error'},false);
            }
            $fileObj.prev('.badge-unsaved').remove();
            if(isNew){
                $fileObj.off('click').removeAttr('data-new').data('new','');
                $fileObj.on('click',function(){
                    openTemplateFile(this,file);
                });
            }
            return App.message({text:r.Info},false);
        },'json');
    };
    save(false);
}
function saveAsTemplateFile(obj,file){
    var $fileObj = $('#templateDirs').find('li.active>a.editable');
    if(file==null) file=$fileObj.data('file');
    var confirmed = false;
    var dir=$('#templateDirs').data('dir');
    var name=$('#templateDirs').data('name');
    var save=function(d,newFile){
        var data = {
            name:name,
            op:'saveAsFileContent',
            dir:dir,
            file:newFile,
            content:$('#console').data('codemirror').getValue(),
        };
        if(confirmed) data.confirmed = 1;
        $.post(window.location.pathname,data,function(r){
            if(r.Code!=1){
                if(r.Code==-109){ // 文件已经存在，确认是否覆盖
                    if(confirm(r.Info)){
                        confirmed = true
                        return save(d,newFile);
                    }
                    return;
                }
                return App.message({text:r.Info,type:'error'},false);
            }
            d.close();
            $fileObj.prev('.badge-unsaved').remove();
            if($('#templateDirs').find('li a.editable[data-file="'+newFile+'"]').length<1){
                $('#templateDirs').children('div').before('<li class="primary-active"><a class="editable" href="javascript:;" data-file="'+newFile+'" onclick="openTemplateFile(this)">\
                <i class="fa fa-file"></i>\
                '+newFile+'\
            </a></li>');
            }
            $('#templateDirs').find('li a.editable[data-file="'+newFile+'"]').trigger('click');
            return App.message({text:r.Info},false);
        },'json');
    };
    showFileNameEditModal(save,file);
}
function removeTemplateFile(obj,file){
    var $fileObj = $('#templateDirs').find('li.active>a.editable');
    if(file==null) file=$fileObj.data('file');
    if(!confirm(App.t('确定要删除文件“%s”吗？',file))) return;
    var dir=$('#templateDirs').data('dir');
    var name=$('#templateDirs').data('name');
    $.post(window.location.pathname,{name:name,op:'removeFile',dir:dir,file:file},function(r){
        if(r.Code!=1){
            return App.message({text:r.Info,type:'error'},false);
        }
        $fileObj.parent('li').remove();
        selectFirstTemplateFile('');
        return App.message({text:r.Info},false);
    },'json');
}
function renameTemplateFile(obj,file){
    var $fileObj = $('#templateDirs').find('li.active>a.editable');
    if(file==null) file=$fileObj.data('file');
    var confirmed = false;
    var dir=$('#templateDirs').data('dir');
    var name=$('#templateDirs').data('name');
    var save=function(d,newFile){
        var data = {
            name:name,
            op:'renameFile',
            dir:dir,
            file:file,
            newFile:newFile,
        };
        if(confirmed) data.confirmed = 1;
        $.post(window.location.pathname,data,function(r){
            if(r.Code!=1){
                if(r.Code==-109){ // 文件已经存在，确认是否覆盖
                    if(confirm(r.Info)){
                        confirmed = true
                        return save(d,newFile);
                    }
                    return;
                }
                return App.message({text:r.Info,type:'error'},false);
            }
            d.close();
            $fileObj.parent('li').remove();
            if($('#templateDirs').find('li a.editable[data-file="'+newFile+'"]').length<1){
                $('#templateDirs').children('div').before('<li class="primary-active"><a class="editable" href="javascript:;" data-file="'+newFile+'" onclick="openTemplateFile(this)">\
                <i class="fa fa-file"></i>\
                '+newFile+'\
            </a></li>');
            }
            $('#templateDirs').find('li a.editable[data-file="'+newFile+'"]').trigger('click');
            return App.message({text:r.Info},false);
        },'json');
    };
    showFileNameEditModal(save,file);
}
function newTemplateFile(obj){
    var save=function(d,newFile){
        if($('#templateDirs').find('li a.editable[data-file="'+newFile+'"]').length>0){
            return App.message({text: App.t('文件“%s”已经存在',newFile), type: 'error'},false);
        }
        var $new=$('<li class="primary-active">\
        <span class="badge-unsaved" title="'+App.t('未保存')+'">*</span>\
        <a class="editable" href="javascript:;" data-file="'+newFile+'" data-new="1">\
            <i class="fa fa-file"></i>\
            '+newFile+'\
        </a></li>');
        $new.find('a').on('click',function(){
            $(this).parent('li').addClass('active').siblings('li.active').removeClass('active');
            $('#console').data('codemirror').setValue('');
        });
        $('#templateDirs').children('div').before($new);
        $new.find('a').trigger('click');
        d.close();
        $('#currentOperateFile').html(newFile);
    };
    showFileNameEditModal(save,'');
}
function showFileNameEditModal(callback,defaultFileName){
    App.editor.dialog({
        title:App.t('设置文件名'),
        message:'<div id="modalFileNameEditBox" class="form-horizontal">\
        <div class="form-group">\
            <label class="col-sm-2 control-label required">'+App.t('文件名')+'</label>\
            <div class="col-sm-10">\
                <input type="text" id="modalInputNewName" value="'+defaultFileName+'" class="form-control" placeholder="'+App.t('文件名')+'" />\
            </div>\
        </div>\
        </div>',
        nl2br:false,
        closeByBackdrop:false,
        onshown: function(d){},
        buttons: [{
          id: 'modifyFileNameDialogBtnSubmit',
          label: App.t('确认'),
          icon: 'fa fa-check',
          cssClass: 'btn-primary mg-r-10',
          action: function(d) {
              var newFileName = $('#modalInputNewName').val();
              if(!newFileName){
                App.message({text:App.t('请输入文件名'),type:'warning'});
                $('#modalInputNewName').focus();
                return;
              }
              callback(d,newFileName);
          }
          },{
          label: App.t('取消'),
          icon: 'fa fa-times',
          cssClass: 'btn-danger',
          action: function(d) {d.close();}
        }]
    });
}
function selectFirstTemplateFile(defaultValue){
    var firstEditable=$('#templateDirs').find('a.editable:first');
    if(firstEditable.length>0) openTemplateFile(firstEditable[0]);
    else if(defaultValue!==null) $('#console').data('codemirror').setValue(defaultValue);
}
$(function(){
$('#pcont').addClass('code-cont');
//$('#pcont').html($('#pcont>.cl-mcont>.main-app'));
$('#pcont>.cl-mcont').replaceWith($('#pcont>.cl-mcont>.main-app'));
var code = '';
App.editor.codemirror('#console', {
    lineNumbers: true,
    theme: 'ambiance',
    value: code,
    mode:  "text/html"
},'html');
setTimeout(function(){
    selectFirstTemplateFile();
    var editor=$('#console').data('codemirror');
    editor.on('change',function(){
        if($('#console').data('onceignorechange')){
            $('#console').data('onceignorechange',false);
            return;
        }
        var $fileObj=$('#templateDirs').find('li.active>a.editable');
        if($fileObj.prev('.badge-unsaved').length<1) {
            var badge='<span class="badge-unsaved" title="'+App.t('未保存')+'">*</span>';
            if($fileObj.prev('.badge-embed').length>0){
                $fileObj.prev('.badge-embed').replaceWith(badge);
            }else{
                $fileObj.before(badge);
            }
        }
    });
},200);
});