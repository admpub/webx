(function(){
function setContainer(callback){
  var $container=$('#content').parent();
  callback.apply($container);
  $('#advertForm').find('textarea[name$="[content]"]').each(function(){
    if(this.id=='content') return;
    var $container=$(this).parent();
    callback.apply($container);
  });
}
$(function(){
  var helpBlock=$('#content-help-block-tpl').html();
  setContainer(function(){
    this.after(helpBlock);
  });
  $('#contype').on('change',function(){
    var cType=$(this).val();
    var callback;
    switch(cType){
      case 'image':
        callback = function(){
          if(this.next('.helper-buttons').length>0)this.next('.helper-buttons').remove();
          var uploadBtn='<input type="file" accept="image/*" data-toggle="uploadPreviewer" data-upload-url="{{$.UploadURL $subdir $urlValues}}" />';
          var browseBtn='<button type="button" class="btn btn-info" data-toggle="finder" data-file-type="'+cType+'"><i class="fa fa-folder"></i> '+App.t('浏览')+'</button>';
          this.after('<div class="helper-buttons">'+browseBtn+'</div>');
          var that=this;
          App.editor.fileInput(that.next('.helper-buttons'),null,function(fileURL){
            that.find('textarea').text(fileURL);
          });
        };
        break;
      case 'video':
        break;
      case 'audio':
        break;
      default:
        break;
    }
    if(callback) setContainer(callback);
  });
});
})();