(function(){
$(function(){
    $('#offlinePayMethod').on('change',function(){
        var v=this.value,$form=$('#offlinePayForm');
        switch(v.split('.')[0]){
            case 'bank':
                $form.find('.method-type-bank.hidden').removeClass('hidden');
                $form.find('.method-type-ebank:not(.hidden)').addClass('hidden');
                break;
            case 'ebank':
                $form.find('.method-type-ebank.hidden').removeClass('hidden');
                $form.find('.method-type-bank:not(.hidden)').addClass('hidden');
                break;
        }
    }).trigger('change');
});
})();