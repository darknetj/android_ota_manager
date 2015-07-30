var $               = require('jquery'),
    _               = require('backbone'),
    Backbone        = require('backbone-crossdomain'),
    NavView         = require('./views/nav.view');
    ReleaseListView = require('./views/release_list.view');

Backbone.$ = $

module.exports = Backbone.Router.extend({
  routes: {
    '': 'releases',
    'releases': 'releases'
  },

  initialize: function() {
    Backbone.history.start();
  },

  releases: function() {
    console.log('releases');
    $(function() {
        new NavView();
        new ReleaseListView();
    });
  }
});
