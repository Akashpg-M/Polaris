export function filterLoadedTwins<T extends { device_id: string; device: { display_name: string } }>(items: T[], search?: string) {
  const value = search?.trim().toLowerCase();
  return value
    ? items.filter(item => item.device_id.toLowerCase().includes(value) || item.device.display_name.toLowerCase().includes(value))
    : items;
}
