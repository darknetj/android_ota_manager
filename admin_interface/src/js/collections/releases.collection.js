var $        = require('jquery'),
    _               = require('backbone'),
    Backbone        = require('backbone-crossdomain'),
    Release  = require("../models/release.model");

module.exports = Backbone.Collection.extend({

    model: Release,

    urlRoot: 'http://localhost:8080',

    url: function() {
        return this.urlRoot + "/releases.json";
    },

    releases: function() {
        return this.toJSON()[0].result;
    },

    // Filter down the list of all todo items that are finished.
    published: function() {
        return this.filter(function(release){ return release.get('published'); });
    },

    // Filter down the list to only todo items that are still not finished.
    unpublished: function() {
        return this.without.apply(this, this.published());
    }
});
