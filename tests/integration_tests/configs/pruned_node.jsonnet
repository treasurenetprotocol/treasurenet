local config = import 'default.jsonnet';

config {
  'treasurenet_5005-1'+: {
    'app-config'+: {
      pruning: 'everything',
      'state-sync'+: {
        'snapshot-interval': 0,
      },
    },
  },
}
