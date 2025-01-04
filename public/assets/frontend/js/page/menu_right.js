$(function() {
    'use strict';
    if($('#showBbMenu').length<1) return;
    // show/hide menu for mobile only
    $('#showBbMenu').on('click', function(){
      if($('body').hasClass('show-bb-menu')) {
        $('body').removeClass('show-bb-menu');
      } else {
        $('body').addClass('show-bb-menu');
      }
      return false;
    });
  
  
    if($('.menu-right').length>0) {
      var cls='';
      if($('.menu-right').children('ul.nav').prev('.pos-absolute-top').length<1) cls=' no-pos-top';
      // add rightpanel for menu wrapper in mobile
      $('body').append('<div class="bb-rightpanel'+cls+'"></div>');
      // function to move menu-right from .headpanel to .rightpanel vice versa
      // fires only when in mobile
      function moveMenuToRight() {
        if($('.menu-right').length<1) return;
        if($('#showBbMenu').css('display') === 'block') {
          $('.menu-right').appendTo('.bb-rightpanel');
        } else {
          var parent=$('.menu-right').data('parent')||'.bb-headpanel .container';
          $('.menu-right').appendTo(parent);
        }
      }
      // calls a function to move .menu-right from .headpanel to .rightpanel
      moveMenuToRight();
      $(window).resize(function(){
        moveMenuToRight();
      });
      // hover dropdown for mega menu
      $('.menu-right .dropdown').hover(function () {
        // Handler in
        if (!$('#showBbMenu').is(':visible')) $(this).addClass('show');
      }, function () {
        // Handler out
        if (!$('#showBbMenu').is(':visible')) $(this).removeClass('show');
      });
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
  });