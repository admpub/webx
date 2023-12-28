$(function(){
    $.get(BACKEND_URL+'/user/message/unread_count',{},function(r){
      if(r.Code!=1) return App.message({title:App.i18n.SYS_INFO,text:r.Info,type:'error'});
      if(r.Data.user>0){
        $('#user-msg-count-label').removeClass('hidden').text(r.Data.user);
        $('#msg-inbox-count-label').removeClass('hidden').text(r.Data.user);
      }
      if(r.Data.system>0){
        $('#system-msg-count-label').removeClass('hidden').text(r.Data.system);
        $('#msg-sys-count-label').removeClass('hidden').text(r.Data.system);
      }
    },'json');
});