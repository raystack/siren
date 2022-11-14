/**
 * Creating a sidebar enables you to:
 - create an ordered group of docs
 - render a sidebar for each doc of that group
 - provide next/previous navigation

 The sidebars can be generated from the filesystem, or explicitly defined here.

 Create as many sidebars as you want.
 */

module.exports = {
  // By default, Docusaurus generates a sidebar from the docs folder structure
  // docsSidebar: [{type: 'autogenerated', dirName: '.'}],

  docsSidebar: [
    'introduction',
    "use_cases",
    'installation',
    {
      type: "category",
      label: "Tour",
      items: [
        "tour/introduction",
        "tour/setup_server",
        "tour/1sending_notifications_overview",
        "tour/2alerting_rules_subscriptions_overview",
      ],
    },
    {
      type: "category",
      label: "Concepts",
      items: [
        "concepts/overview",
        "concepts/plugin",
        "concepts/notification",
        "concepts/glossary",
      ],
    },
    {
      type: "category",
      label: "Guides",
      items: [
        "guides/overview",
        "guides/deployment",
        "guides/provider_and_namespace",
        "guides/receiver",
        "guides/subscription",
        "guides/rule",
        "guides/template",
        "guides/alert_history",
        "guides/notification",
        "guides/workers",
        "guides/job",
      ],
    },
    {
      type: "category",
      label: "Providers",
      items: [
        "providers/cortexmetrics",
      ],
    },
    {
      type: "category",
      label: "Receivers",
      items: [
        "receivers/slack",
        "receivers/pagerduty",
        "receivers/http",
        "receivers/file",
      ],
    },
    {
      type: "category",
      label: "Reference",
      items: [
        "reference/api",
        "reference/server_configuration",
        "reference/client_configuration",
        "reference/cli",
      ],
    },
    {
      type: "category",
      label: "Extend",
      items: [
        "extend/adding_new_provider",
        "extend/adding_new_receiver"
      ],
    },
    {
      type: "category",
      label: "Contribute",
      items: [
        "contribute/contribution",
        "contribute/release"],
    },
  ],
};
