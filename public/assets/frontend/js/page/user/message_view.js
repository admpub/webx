$(function(){
    App.editor.markdownToHTML($('#message-content'));
    App.editor.markdownToHTML($('#reply-list-box'));
    $('#input-reply-content').on('keyup',function(e){
        if(e.keyCode!=13)return;
        var data={replyId: $(this).data('reply-id'),content: $(this).val()};
        var ajaxOptions={
            url: BASE_URL+'/user/message/send',
            type: 'POST',
            data: data,
            dataType: 'json',
            success: function(r){
                onAjaxRespond(null,r,ajaxOptions,function(r){
                    window.setTimeout(function(){window.location.reload();},1000);
                });
            }
        };
        $.ajax(ajaxOptions);
    });
});