(function(){
$(function(){
    $('#offlinePayMethod').on('change',function(){
        var v=this.value,$form=$('#offlinePayForm');
        switch(v.split('.')[0]){
            case 'bank':
                $form.find('.method-type-bank.hidden').removeClass('hidden');
                $form.find('.method-type-ebank:not(.hidden)').addClass('hidden');
                $('#offlinePayTransactionNo').prop('required',false);
                $('#offlinePayTransactionNo').parent().prev('label').removeClass('required');
                break;
            case 'ebank':
                $form.find('.method-type-ebank.hidden').removeClass('hidden');
                $form.find('.method-type-bank:not(.hidden)').addClass('hidden');
                $('#offlinePayTransactionNo').prop('required',true);
                $('#offlinePayTransactionNo').parent().prev('label').addClass('required');
                break;
        }
    }).trigger('change');
});
})();