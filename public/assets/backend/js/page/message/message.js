$(function(){
    $.get(BACKEND_URL+'/user/message/unread_count',{},function(r){
      if(r.Code!=1) return App.message({title:App.i18n.SYS_INFO,text:r.Info,type:'error'});
      var hasNew = false;
      if(r.Data.user>0){
        $('#user-msg-count-label').removeClass('hidden').text(r.Data.user);
        $('#msg-inbox-count-label').removeClass('hidden').text(r.Data.user);
        hasNew = true;
      }
      if(r.Data.system>0){
        $('#system-msg-count-label').removeClass('hidden').text(r.Data.system);
        $('#msg-sys-count-label').removeClass('hidden').text(r.Data.system);
        hasNew = true;
      }
      if(hasNew){
        $('#square-buttons-hasnew').removeClass('hidden');
      }
    },'json');
});