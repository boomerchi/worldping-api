define([
  'angular',
  'app',
  'lodash',
  'config',
  'components/panelmeta',
],
function (angular, app, _, config, PanelMeta) {
  'use strict';

  var module = angular.module('grafana.panels.dashlist', []);
  app.useModule(module);

  module.directive('grafanaPanelDashlist', function() {
    return {
      controller: 'DashListPanelCtrl',
      templateUrl: 'app/panels/dashlist/module.html',
    };
  });

  module.controller('DashListPanelCtrl', function($scope, panelSrv, backendSrv) {

    $scope.panelMeta = new PanelMeta({
      panelName: 'Dash list',
      editIcon:  "fa fa-star",
      fullscreen: true,
    });

    $scope.panelMeta.addEditorTab('Options', 'app/panels/dashlist/editor.html');

    var defaults = {
      mode: 'starred',
      query: '',
      limit: 10,
      tag: '',
    };

    $scope.modes = ['starred', 'search'];

    _.defaults($scope.panel, defaults);

    $scope.dashList = [];

    $scope.init = function() {
      panelSrv.init($scope);

      if ($scope.isNewPanel()) {
        $scope.panel.title = "Starred Dashboards";
      }
    };

    $scope.refreshData = function() {
      var params = {
        limit: $scope.panel.limit
      };

      if ($scope.panel.mode === 'starred') {
        params.starred = "true";
      } else {
        params.query = $scope.panel.query;
        params.tag = $scope.panel.tag;
      }

      return backendSrv.search(params).then(function(result) {
        $scope.dashList = result.dashboards;
      });
    };

    $scope.init();
  });
});
