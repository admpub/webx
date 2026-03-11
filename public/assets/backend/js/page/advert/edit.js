(function(){
var formE='#advertForm';
function setContainer(callback){
  var $container=$('#content').parent();
  callback.apply($container);
  $(formE).find('textarea[name$="[content]"]').each(function(){
    if(this.id=='content') return;
    var $container=$(this).parent();
    callback.apply($container);
  });
}
function setHelperButtons(cType){
  if(this.next('.helper-buttons').length>0)this.next('.helper-buttons').remove();
  var uploadBtn='<input type="file" accept="'+cType+'/*" data-toggle="uploadPreviewer" data-upload-url="'+$(formE).data('upload-url')+'" />';
  var browseBtn='<button type="button" class="btn btn-info" data-toggle="finder" data-file-type="'+cType+'" data-url-query="subdir=advert"><i class="fa fa-folder"></i> '+App.t('浏览')+'</button>';
  var previewContent='', value=this.find('textarea').text();
  switch(cType){
    case 'image': previewContent='<img class="preview" src="'+value+'" style="max-width:100%" />'; break;
    case 'video': previewContent='<video class="preview" controls="controls" src="'+value+'" style="max-width:100%"></video>'; break;
    case 'audio': previewContent='<audio class="preview" controls="controls" src="'+value+'" style="max-width:100%"></audio>'; break;
  }
  var previewBox='<div>'+previewContent+'</div>';
  this.after('<div class="helper-buttons">'+uploadBtn+browseBtn+previewBox+'</div>');
  var that=this, helperButtons=that.next('.helper-buttons');
  App.editor.fileInput(helperButtons,null,function(fileURL){
    that.find('textarea').text(fileURL);
    helperButtons.find('.preview').attr('src',fileURL);
  });
}
$(function(){
  var helpBlock=$('#content-help-block-tpl').html();
  setContainer(function(){
    this.after(helpBlock);
    var that=this;
    that.find('textarea').on('change',function(){
      var helperButtons=that.next('.helper-buttons');
      if(helperButtons.length<1)return;
      helperButtons.find('.preview').attr('src',$(this).val());
    });
  });
  $('#contype').on('change',function(){
    var cType=$(this).val();
    var callback;
    switch(cType){
      case 'image':case 'video':case 'audio':callback = function(){setHelperButtons.call(this,cType)};break;
      default:callback = function(){if(this.next('.helper-buttons').length>0)this.next('.helper-buttons').remove();};break;
    }
    if(callback) setContainer(callback);
  }).trigger('change');
});
})();