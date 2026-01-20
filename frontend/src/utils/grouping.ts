import { ContainerInfo } from '../services/api';

export interface GroupedContainers {
  [key: string]: ContainerInfo[];
}

export function groupContainers(containers: ContainerInfo[], groupBy: string): GroupedContainers {
  if (groupBy === 'none') {
    return { 'All Containers': containers };
  }

  const grouped: GroupedContainers = {};

  containers.forEach((container) => {
    let groupKey = 'Uncategorized';

    switch (groupBy) {
      case 'project':
        groupKey = container.labels?.['com.docker.compose.project'] || 
                   container.labels?.['project'] || 
                   container.labels?.['com.project'] ||
                   'Uncategorized';
        break;
      case 'environment':
        groupKey = container.labels?.['com.docker.compose.project.working_dir'] ||
                   container.labels?.['environment'] ||
                   container.labels?.['env'] ||
                   container.labels?.['com.environment'] ||
                   'Uncategorized';
        break;
      case 'image':
        // Extract image name without tag
        const imageName = container.image.split(':')[0];
        groupKey = imageName || 'Unknown';
        break;
      case 'state':
        groupKey = container.state || 'unknown';
        break;
      default:
        groupKey = 'All Containers';
    }

    if (!grouped[groupKey]) {
      grouped[groupKey] = [];
    }
    grouped[groupKey].push(container);
  });

  return grouped;
}


