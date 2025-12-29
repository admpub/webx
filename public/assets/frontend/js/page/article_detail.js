$(function(){
    $('#declare-like').on('click',function(){
        var id=$(this).data('article-id'),me=$(this);
        $.post(BASE_URL+'/article/like',{id:id},function(r){
            if(r.Code<1){
                App.message({text:r.Info,type:"error"});
                return;
            }
            me.removeClass('tx-gray-light').addClass('tx-success');
            var numElem=$('#article-likes');
            if(numElem.length>0){
                numElem.text(numElem.text()*1+1);
            }
        });
    });
    $('#declare-hate').on('click',function(){
        var id=$(this).data('article-id'),me=$(this);
        $.post(BASE_URL+'/article/hate',{id:id},function(r){
            if(r.Code<1){
                App.message({text:r.Info,type:"error"});
                return;
            }
            me.removeClass('tx-gray-light').addClass('tx-danger');
            var numElem=$('#article-hates');
            if(numElem.length>0){
                numElem.text(numElem.text()*1+1);
            }
        });
    });
    $('#collect-article').on('click',function(){
        var id=$(this).data('article-id'),me=$(this);
        $.post(BASE_URL+'/article/collect',{id:id},function(r){
            if(r.Code<1){
                App.message({text:r.Info,type:"error"});
                return;
            }
            if(r.Data && r.Data.cancel){
                me.removeClass('tx-success').addClass('tx-gray-light');
                me.children('span').text(App.t('收藏'));
            }else{
                me.removeClass('tx-gray-light').addClass('tx-success');
                me.children('span').text(App.t('已收藏'));
            }
        });
    });
    $('.content').on('click','img',function(){
        $(this).toggleClass('img-full');
    });
    $('body').on('click','.comment-content img',function(){
        $(this).toggleClass('img-full');
    });
    $('#article-content').on('click','[data-article-payment]',function(){
        var id = $(this).data('article-payment');
        var price = $(this).data('article-price');
        $.post(BASE_URL+'/article/pay/'+id,function(r){
            if(r.Code!=1) {
                switch(r.Code){
                    case -1: // 未登录
                        window.setTimeout(function(){
                            window.location.href=BASE_URL+'/sign_in?next='+encodeURIComponent(window.location.href);
                        },2000);
                        break;
                    case -5: // 余额不足
                        r.Info=App.t(r.Info)+', '+App.t('前往充值...');
                        window.setTimeout(function(){
                            window.location.href=BASE_URL+'/user/wallet/recharge?amount='+price+'&next='+encodeURIComponent(window.location.href);
                        },2000);
                        break;
                }
                return App.message({text:r.Info,type:'error'});
            }
            App.message({text:App.t('购买成功，重新加载中...'),type:'success'});
            window.setTimeout(function(){
                window.location.reload();
            },2000);
        },'json');
    });
});