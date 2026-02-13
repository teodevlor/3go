/**
 * Tài liệu API cho Go Structure Project
 * Tổ chức theo các module chính: Website System, App Driver, App User
 * @type {import('@docusaurus/plugin-content-docs').SidebarsConfig}
 */
const sidebars = {
  apiSidebar: [
    'intro',
    {
      type: 'category',
      label: 'Website System',
      collapsed: false,
      link: {
        type: 'generated-index',
        title: 'Website System API',
        description: 'API quản lý hệ thống website: Zones (khu vực/vùng), Sidebars (menu động).',
        slug: '/website-system',
      },
      items: [
        'website-system/zones',
        'website-system/sidebars',
      ],
    },
    {
      type: 'category',
      label: 'App User',
      collapsed: false,
      link: {
        type: 'generated-index',
        title: 'App User API',
        description: 'API cho người dùng: Xác thực, quản lý profile, mật khẩu và OTP.',
        slug: '/app-user',
      },
      items: [
        'app-user/authentication',
        'app-user/profile',
        'app-user/password',
        'app-user/otp',
      ],
    },
  ],
};

module.exports = sidebars;
