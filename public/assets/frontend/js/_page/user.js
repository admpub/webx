$(function() {
  // add rightpanel for menu wrapper in mobile
  $('body').append('<div class="bb-rightpanel"></div>');

  // show/hide menu for mobile only
  $('#showBbMenu').on('click', function(){
    if($('body').hasClass('show-bb-menu')) {
      $('body').removeClass('show-bb-menu');
    } else {
      $('body').addClass('show-bb-menu');
    }
    return false;
  });

  // function to move menu-right from .headpanel to .rightpanel vice versa
  // fires only when in mobile
  function moveMenuToRight() {
    if($('#showBbMenu').css('display') === 'block') {
      $('.menu-right').appendTo('.bb-rightpanel');
    } else {
      $('.menu-right').appendTo('.bb-headpanel .container');
    }
  }

  // calls a function to move .menu-right from .headpanel to .rightpanel
  moveMenuToRight();
  $(window).resize(function(){
    moveMenuToRight();
  });

  // Select2 without the search
  if($.fn.select2) {
    $('.select2').select2({
      minimumResultsForSearch: Infinity
    });

    // Select2 by showing the search
    $('.select2-show-search').select2({
      minimumResultsForSearch: ''
    });

    // Select2 with tagging support
    $('.select2-tag').select2({
      tags: true,
      tokenSeparators: [',', ' ']
    });
  }

  // Datepicker
  if($.fn.datepicker) {
    $('.form-control-datepicker').datepicker();
  }

  // function to move menu-right from .headpanel to .rightpanel vice versa
  // fires only when in mobile
  // only fires to templates/components-*.html files
  if($('.component-menu').length) {
    function moveComponentsMenuToRight() {
      if($('#showBbMenu').css('display') === 'block') {
        $('.component-menu').appendTo('.component-menu-sidebar-wrapper');
      } else {
        $('.component-menu').appendTo('.component-menu-wrapper');
      }
    }

    // calls a function to move .menu-right from .headpanel to .rightpanel
    moveComponentsMenuToRight();
    $(window).resize(function(){
      moveComponentsMenuToRight();
    });
  }

    // show/hide search bar in smaller screen devices only
    $('#searchmenu').on('click', function(e) {
      if($('body').hasClass('show-search')) {
        $('body').removeClass('show-search');
      } else {
        $('body').addClass('show-search');
      }
    });

    $('.left-menu').on('click', function(){
      if($('body').hasClass('show-left-menu')) {
        $('body').removeClass('show-left-menu');
      } else {
        $('body').addClass('show-left-menu');
      }
    });
    $.get(BASE_URL+'/user/message/unread_count',{},function(r){
      if(r.Code!=1) return showMsg({text:r.Info,type:'error'});
      if(r.Data.user>0)$('#user-msg-count-square').removeClass('d-none');
      if(r.Data.system>0)$('#system-msg-count-square').removeClass('d-none');
    },'json');
});