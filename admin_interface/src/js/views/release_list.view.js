var $                  = require('jquery'),
    _                  = require('backbone'),
    Backbone           = require('backbone-crossdomain'),
    ReleasesCollection = require('../collections/releases.collection')
    template           = require('../templates/release_list.hbs');

Backbone.$ = window.$

module.exports = Backbone.View.extend({
    el: $('#main'),
    template: template,

    initialize: function () {
        console.log("releases initialize");
        var releases = this;
        releases.collection = new ReleasesCollection();
        releases.collection.on("reset", this.render, this);
        releases.collection.fetch({
            success: function () {
                releases.render();
            }
        });
    },

    render: function () {
        this.$el.html(template({ releases: this.collection.releases() }));
        return this;
    }
});
